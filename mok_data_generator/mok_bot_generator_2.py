import discord
from discord.ext import tasks, commands
import os
import random
import asyncio
import datetime
from openai import AzureOpenAI
import logging

# --- Logging Setup ---
logging.basicConfig(level=logging.INFO, format='%(asctime)s:%(levelname)s:%(name)s: %(message)s')
logger = logging.getLogger('discord.technical_bot')

# --- Essential Credentials Loading ---
BOT_TOKEN = os.getenv('BOT_TOKEN_2')
AZURE_ENDPOINT = os.getenv('AZURE_OPENAI_ENDPOINT')
AZURE_KEY = os.getenv('AZURE_OPENAI_API_KEY')
AZURE_DEPLOYMENT_NAME = os.getenv('AZURE_OPENAI_DEPLOYMENT_NAME')

# --- Bot Behavior Configuration (Hardcoded) ---
# These values are now set directly in the code
CONVERSATION_START_INTERVAL_HOURS = 6.0
TECHNICAL_RESPOND_CHANCE = 1
MAX_TECHNICAL_RESPONSE_TOKENS = 150
MAX_TECHNICAL_STARTER_TOKENS = 3000
POST_STARTER_IF_INACTIVE_HOURS = 3.0
# TARGET_CHANNEL_IDS is no longer used

# --- Validate Essential Configuration ---
# Only check essential credentials now
if not all([BOT_TOKEN, AZURE_ENDPOINT, AZURE_KEY, AZURE_DEPLOYMENT_NAME]):
    raise ValueError("One or more required environment variables (BOT_TOKEN, AZURE_OPENAI_ENDPOINT, AZURE_OPENAI_API_KEY, AZURE_OPENAI_DEPLOYMENT_NAME) are missing in .env_technical.")

logger.info(f"Conversation Start Interval: {CONVERSATION_START_INTERVAL_HOURS} hours")
logger.info(f"Technical Respond Chance: {TECHNICAL_RESPOND_CHANCE}")
logger.info(f"Post Starter If Inactive Hours: {POST_STARTER_IF_INACTIVE_HOURS}")

# --- Azure OpenAI Client Setup ---
# (Remains the same)
AZURE_API_VERSION = "2024-02-01"
azure_client = None
try:
    azure_client = AzureOpenAI(
        api_key=AZURE_KEY,
        api_version=AZURE_API_VERSION,
        azure_endpoint=AZURE_ENDPOINT,
    )
    logger.info("Azure OpenAI client initialized successfully.")
except Exception as e:
    logger.error(f"Error initializing Azure OpenAI client: {e}")

# --- Bot Setup ---
# (Remains the same)
intents = discord.Intents.default()
intents.messages = True
intents.message_content = True
intents.guilds = True # Still needed to iterate guilds

bot = commands.Bot(command_prefix="tech!", intents=intents) # Prefix might not be used

# --- Helper Functions (generate_technical_content) ---
# (Remains the same as the previous version)
async def generate_technical_content(mode: str, channel_name: str = "", history: list = [], user_message: discord.Message = None) -> str | None:
    # [ Keep the exact same function code from the previous version here ]
    # ... (ensure the full function is copied here) ...
    if not azure_client:
        logger.error("Azure OpenAI client is not available.")
        return None

    messages_for_api = []
    max_tokens = MAX_TECHNICAL_RESPONSE_TOKENS
    temperature = 0.5 # Lower temp for more factual/focused technical responses

    if mode == "starter":
        max_tokens = MAX_TECHNICAL_STARTER_TOKENS
        temperature = 0.7 # Slightly higher temp for more varied starters
        system_prompt = """You are an AI assistant tasked with initiating engaging technical discussions on Discord.
Generate a thought-provoking, open-ended technical question or a brief statement to spark conversation.
Topics can include programming languages, software architecture, AI/ML, cybersecurity, hardware, cloud computing, data science, or interesting scientific concepts.
Keep it concise and focused. Avoid greetings or introductions. Generate only the question or statement itself."""
        messages_for_api.append({"role": "system", "content": system_prompt})
        messages_for_api.append({"role": "user", "content": "Generate a technical conversation starter."}) # Simple trigger

    elif mode == "response" and user_message and channel_name:
        system_prompt = f"""You are a highly knowledgeable technical expert participating in a Discord chat in the channel #{channel_name}.
Your purpose is to provide accurate, detailed, and informative technical responses. Maintain a professional, objective, and helpful tone.
When appropriate, elaborate on concepts, provide examples, explain your reasoning, or reference relevant principles. Avoid slang and overly casual language.
You are responding to '{user_message.author.display_name}'. Analyze their message within the context of the recent technical discussion history provided below.
Generate a clear and relevant technical contribution to the conversation."""
        messages_for_api.append({"role": "system", "content": system_prompt})

        # Add history (older first)
        for msg in reversed(history):
            role = "assistant" if msg.author == bot.user else "user"
            sanitized_name = "".join(c for c in msg.author.display_name.replace(" ", "_") if c.isalnum() or c in ['_', '-'])[:64] or ("user" if role == "user" else "assistant")
            messages_for_api.append({"role": role, "name": sanitized_name, "content": msg.clean_content})
        # Triggering message is already included in the history list passed to this function

    else:
        logger.error(f"Invalid mode or missing arguments for generate_technical_content: mode={mode}")
        return None

    # --- Call Azure OpenAI API ---
    try:
        logger.info(f"Calling Azure OpenAI (Mode: {mode}, Deployment: {AZURE_DEPLOYMENT_NAME})")
        # logger.debug(f"Messages for API: {messages_for_api}") # Debug

        response = await asyncio.to_thread(
            azure_client.chat.completions.create,
            model=AZURE_DEPLOYMENT_NAME,
            messages=messages_for_api,
            max_tokens=max_tokens,
            temperature=temperature,
            top_p=1.0,
            n=1,
            stop=None
        )

        ai_content = response.choices[0].message.content.strip()
        logger.info(f"AI Response Raw: '{ai_content}'")

        # Basic cleanup
        if ai_content.startswith('"') and ai_content.endswith('"'):
            ai_content = ai_content[1:-1].strip()
        if mode == "starter" and (":" in ai_content or ai_content.startswith("Sure,") or ai_content.startswith("Here's")) :
            # Attempt to clean up introductory phrases sometimes added by models
            parts = ai_content.split(':')
            if len(parts) > 1: ai_content = parts[1].strip()
            ai_content = ai_content.removeprefix("Sure, ").removeprefix("Here's ")
            ai_content = ai_content.strip()


        if not ai_content:
            logger.warning("AI returned an empty response.")
            return None

        return ai_content

    except Exception as e:
        logger.exception(f"Error calling Azure OpenAI: {e}") # Log full traceback
        return None


# --- Tasks ---

@tasks.loop(hours=CONVERSATION_START_INTERVAL_HOURS)
async def start_conversation_task():
    """Periodically tries to start a technical conversation in a random accessible text channel."""
    if not azure_client:
        logger.warning("start_conversation_task: Skipping because Azure client is not available.")
        return

    # --- Find all suitable channels across all guilds ---
    eligible_channels = []
    for guild in bot.guilds:
        logger.debug(f"Checking guild: {guild.name} ({guild.id})")
        for channel in guild.text_channels:
            # Check if the bot has permissions to send messages in this channel
            perms = channel.permissions_for(guild.me) # Get bot's permissions in this specific channel
            if perms.send_messages:
                # Additional check for read history if inactivity check is enabled
                eligible_channels.append(channel)
            else:
                logger.debug(f"Skipping channel {channel.name} ({channel.id}) in {guild.name}: Missing Send Messages permission.")


    if not eligible_channels:
        logger.warning("start_conversation_task: No eligible text channels found across all guilds where the bot can send messages.")
        return

    # --- Choose one random channel from the eligible list ---
    target_channel = random.choice(eligible_channels)
    guild = target_channel.guild # Get the guild for logging context

    logger.info(f"start_conversation_task: Selected target channel: {target_channel.name} ({target_channel.id}) in guild {guild.name} ({guild.id})")

    # --- Optional: Check for recent activity in the chosen channel ---
    if POST_STARTER_IF_INACTIVE_HOURS > 0:
        try:
            # Fetch the very last message; requires read_message_history
            last_message = await target_channel.fetch_message(target_channel.last_message_id) if target_channel.last_message_id else None

            if last_message:
                now_utc = datetime.datetime.now(datetime.timezone.utc)
                time_since_last_message = now_utc - last_message.created_at
            else:
                logger.info(f"start_conversation_task: Channel {target_channel.name} appears empty or last message inaccessible; proceeding.")

        except discord.errors.Forbidden:
            logger.warning(f"start_conversation_task: No permission to read history/last message in {target_channel.name}. Cannot check inactivity.")
            # Proceed without inactivity check if permissions are missing
        except Exception as e:
            logger.exception(f"start_conversation_task: Error checking last message in {target_channel.name}: {e}")
            # Proceed cautiously if check fails

    # --- Generate and Send Starter ---
    starter_prompt = await generate_technical_content(mode="starter")

    if starter_prompt:
        try:
            await target_channel.send(starter_prompt)
            logger.info(f"Successfully sent conversation starter to {target_channel.name}: '{starter_prompt}'")
        except discord.errors.Forbidden:
            # This check *should* be redundant due to earlier filtering, but good as a safeguard
            logger.error(f"start_conversation_task: Bot lost/lacks send permissions in channel {target_channel.name} ({target_channel.id}).")
        except Exception as e:
            logger.exception(f"start_conversation_task: Failed to send message to {target_channel.name} ({target_channel.id}): {e}")
    else:
        logger.error("start_conversation_task: Failed to generate starter prompt from AI.")


@start_conversation_task.before_loop
async def before_start_conversation_task():
    # (Remains the same)
    logger.info("Waiting for bot to be ready before starting conversation task...")
    await bot.wait_until_ready()
    logger.info("Bot ready. Conversation starting task loop begins.")

# --- Bot Events ---

@bot.event
async def on_ready():
    # (Remains the same)
    logger.info(f'Logged in as {bot.user.name} ({bot.user.id})')
    logger.info(f"Bot will operate in ALL accessible text channels.")
    if not azure_client:
        logger.warning("Bot started but Azure OpenAI client failed to initialize. AI features disabled.")
    else:
        logger.info("Starting background task: start_conversation_task")
        start_conversation_task.start()
    await bot.change_presence(activity=discord.Activity(type=discord.ActivityType.listening, name="technical discussions"))
    logger.info('------ Technical Bot Ready ------')

@bot.event
async def on_message(message):
    """Handles incoming messages in any accessible channel."""
    # 1. Ignore self, bots, DMs
    if message.author == bot.user:
        return

    # 2. REMOVED Channel ID Check - Bot responds in any channel subject to chance

    # 3. Check Azure client
    if not azure_client:
        return

    # 4. Random chance to respond
    if random.random() > TECHNICAL_RESPOND_CHANCE:
        return

    # --- Proceed with generating a technical response ---
    # Log channel name for context
    logger.info(f"\n[{datetime.datetime.now()}] Considering technical response in '{message.channel.name}' ({message.channel.id}) in guild '{message.guild.name}' to {message.author.name}: '{message.content}'")

    # Check permissions needed for responding (Read History, Send Messages)
    perms = message.channel.permissions_for(message.guild.me)
    if not perms.read_message_history or not perms.send_messages:
        logger.warning(f"on_message: Missing Read History or Send Messages permission in {message.channel.name}. Cannot respond.")
        return

    message_history = []
    try:
        async for msg in message.channel.history(limit=8, before=message):
            message_history.append(msg)
        message_history.insert(0, message) # Add triggering message
        logger.info(f"Fetched {len(message_history)-1} messages for technical history in {message.channel.name}.")
    # Forbidden error less likely now due to check above, but keep for safety
    except discord.errors.Forbidden:
        logger.error(f"on_message: Bot lost/lacks permissions to read message history in channel '{message.channel.name}'")
        return
    except Exception as e:
        logger.exception(f"on_message: Error fetching message history in {message.channel.name}: {e}")
        message_history = [message] # Fallback
        logger.warning("Warning: Proceeding with AI generation using only the triggering message.")

    # Generate response using AI with the 'response' mode prompt
    ai_response = await generate_technical_content(
        mode="response",
        channel_name=message.channel.name,
        history=message_history,
        user_message=message
    )

    if not ai_response:
        logger.warning("AI technical response generation failed or returned empty.")
        return

    # Simulate slight delay and send
    try:
        delay = random.uniform(0.8, 2.5)
        await asyncio.sleep(delay)

        await message.channel.send(ai_response)
        logger.info(f"Sent technical response in '{message.channel.name}' to {message.author.name}: '{ai_response}'")

    # Forbidden error less likely now due to check above, but keep for safety
    except discord.errors.Forbidden:
        logger.error(f"on_message: Bot lost/lacks permissions to send message in channel '{message.channel.name}'")
    except Exception as e:
        logger.exception(f"on_message: An unexpected error occurred during response sending: {e}")


# --- Run the Bot ---
# (Remains the same - uses BOT_TOKEN)
if __name__ == "__main__":
    if not BOT_TOKEN:
        logger.critical("ERROR: Discord Bot token (BOT_TOKEN) not found in .env_technical file.")
    else:
        try:
            logger.info("Starting Technical Discord bot...")
            bot.run(BOT_TOKEN)
        except discord.errors.LoginFailure:
            logger.critical("ERROR: Failed to log in to Discord. Check BOT_TOKEN in .env_technical.")
        except Exception as e:
            logger.critical(f"An error occurred while running the Discord bot: {e}", exc_info=True)