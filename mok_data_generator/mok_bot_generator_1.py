import discord
import os
import random
import asyncio
import datetime
from openai import AzureOpenAI
import logging # Added logging

# --- Logging Setup ---
logging.basicConfig(level=logging.INFO, format='%(asctime)s:%(levelname)s:%(name)s: %(message)s')
logger = logging.getLogger('discord.always_responder_bot')

# --- Configuration ---
# Load essential credentials from .env file (make sure it exists)
BOT_TOKEN = os.getenv('BOT_TOKEN')
AZURE_ENDPOINT = os.getenv('AZURE_OPENAI_ENDPOINT')
AZURE_KEY = os.getenv('AZURE_OPENAI_API_KEY')
AZURE_DEPLOYMENT_NAME = os.getenv('AZURE_OPENAI_DEPLOYMENT_NAME')

# --- Bot Behavior Tuning (Hardcoded for 100% Response) ---
# CONSIDER_RESPONDING_CHANCE is removed (always responds)
# REACTION_ONLY_CHANCE is removed (always generates text)
# APPEND_EMOJI_CHANCE is removed (no random emojis appended)

# Max number of recent messages to fetch for context
MAX_HISTORY = 6

# Max tokens for the AI text response (Reduced for frequent responses)
MAX_RESPONSE_TOKENS = 300 # Lowered from 350

# COMMON_REACTIONS list is removed (no longer needed)
# APPENDABLE_EMOJIS list is removed (no longer needed)

# --- Validate Configuration ---
if not all([BOT_TOKEN, AZURE_ENDPOINT, AZURE_KEY, AZURE_DEPLOYMENT_NAME]):
    # Use logging instead of raising ValueError immediately if possible
    logger.critical("CRITICAL: One or more required environment variables (BOT_TOKEN, AZURE_OPENAI_ENDPOINT, AZURE_OPENAI_API_KEY, AZURE_OPENAI_DEPLOYMENT_NAME) are missing. Bot cannot start.")
    # You might want to exit gracefully here depending on deployment strategy
    exit("Missing essential environment variables.") # Exit if critical config is missing

# --- Azure OpenAI Client Setup ---
# *** YOU MAY NEED TO ADJUST THIS API VERSION ***
AZURE_API_VERSION = "2024-02-01"
azure_client = None # Initialize as None
try:
    azure_client = AzureOpenAI(
        api_key=AZURE_KEY,
        api_version=AZURE_API_VERSION,
        azure_endpoint=AZURE_ENDPOINT,
    )
    logger.info("Azure OpenAI client initialized successfully.")
except Exception as e:
    logger.error(f"Error initializing Azure OpenAI client: {e}. AI features will be disabled.")
    # Bot can continue running, but AI calls will fail.

# --- Bot Setup ---
intents = discord.Intents.default()
intents.messages = True
intents.message_content = True
intents.guilds = True

discord_client = discord.Client(intents=intents)

# --- Helper Functions ---

async def generate_ai_response(channel_name: str, history: list[discord.Message], user_message: discord.Message) -> str | None:
    """Generates a response using Azure OpenAI, simulating human chat and aiming to keep conversation alive."""
    if not azure_client:
        logger.error("Azure OpenAI client is not initialized. Cannot generate response.")
        return None

    bot_persona_name = discord_client.user.display_name

    # --- Modified System Prompt ---
    system_prompt = f"""You are a highly knowledgeable technical expert participating in a Discord chat in the channel #{channel_name}.
Your purpose is to provide accurate, detailed, and informative technical responses. Maintain a professional, objective, and helpful tone.
When appropriate, elaborate on concepts, provide examples, explain your reasoning, or reference relevant principles. Avoid slang and overly casual language.
You are responding to '{user_message.author.display_name}'. Analyze their message within the context of the recent technical discussion history provided below.
Generate a clear and relevant technical contribution to the conversation."""


    messages_for_api = [{"role": "system", "content": system_prompt}]

    for msg in reversed(history):
        role = "assistant" if msg.author == discord_client.user else "user"
        sanitized_name = "".join(c for c in msg.author.display_name.replace(" ", "_") if c.isalnum() or c in ['_', '-'])[:64]
        if not sanitized_name: sanitized_name = "user" if role == "user" else "assistant"
        messages_for_api.append({"role": role, "name": sanitized_name, "content": msg.clean_content})

    # --- Call Azure OpenAI API ---
    try:
        logger.info(f"Calling Azure OpenAI for channel #{channel_name}")
        # logger.debug(f"Messages sent to API: {messages_for_api}") # Uncomment for debugging prompt

        response = await asyncio.to_thread(
            azure_client.chat.completions.create,
            model=AZURE_DEPLOYMENT_NAME,
            messages=messages_for_api,
            max_tokens=300,
            temperature=1.0,
            top_p=1.0,
            n=1,
            stop=None
        )

        ai_content = response.choices[0].message.content.strip()
        logger.info(f"AI Response Raw: '{ai_content}'")

        # Basic cleanup
        if ai_content.startswith('"') and ai_content.endswith('"'):
            ai_content = ai_content[1:-1].strip()
        if ai_content.lower().startswith(f"{bot_persona_name.lower()}:"):
            ai_content = ai_content[len(bot_persona_name)+1:].strip()

        if not ai_content:
            logger.warning("AI returned an empty response.")
            return None

        # Removed the random emoji append logic

        return ai_content

    except Exception as e:
        # Log the full exception for better debugging
        logger.exception(f"Error calling Azure OpenAI: {e}")
        return None

# --- Bot Events ---

@discord_client.event
async def on_ready():
    """Runs when the bot successfully connects."""
    logger.info(f'Logged in as {discord_client.user.name} ({discord_client.user.id})')
    logger.warning("Bot configured to respond to ALL messages in ALL accessible channels.") # Add warning
    if not azure_client:
        logger.error("Azure OpenAI client failed to initialize. AI responses will NOT work.")

    # Set presence
    try:
        # Simplified presence
        await discord_client.change_presence(activity=discord.Activity(type=discord.ActivityType.listening, name="everything"))
        logger.info("Set presence to 'listening to everything'")
    except Exception as e:
        logger.error(f"Failed to set presence: {e}")
    logger.info('------ Always Responder Bot Ready ------')


@discord_client.event
async def on_message(message):
    """Runs whenever a message is sent. Responds to ALL messages."""
    # 1. Ignore self, bots, DMs
    if message.author == discord_client.user or not message.guild:
        return

    # 2. Check Azure client (essential for responding)
    if not azure_client:
        # Silently return if AI isn't available, maybe log periodically elsewhere if needed
        # logger.debug("Skipping response: Azure client not ready.") # Avoid logging this every message
        return

    # --- Always proceed to generate text response ---
    # Removed the CONSIDER_RESPONDING_CHANCE check
    # Removed the REACTION_ONLY_CHANCE check

    channel_name = message.channel.name
    guild_name = message.guild.name
    logger.info(f"\nProcessing message in #{channel_name} ({guild_name}) from {message.author.name}: '{message.content}'")

    # --- Check permissions before proceeding ---
    perms = message.channel.permissions_for(message.guild.me)
    if not perms.send_messages or not perms.read_message_history:
        logger.warning(f"Missing Send Messages or Read History permission in #{channel_name} ({guild_name}). Cannot respond.")
        return

    # --- Fetch History ---
    message_history = []
    try:
        async for msg in message.channel.history(limit=MAX_HISTORY, before=message):
            message_history.append(msg)
        message_history.insert(0, message)
        logger.info(f"Fetched {len(message_history)-1} messages for history in #{channel_name}.")
    except discord.errors.Forbidden: # Should be caught by check above, but good safety
        logger.error(f"Permission error fetching history in #{channel_name} despite initial check.")
        return
    except Exception as e:
        logger.exception(f"Error fetching message history in #{channel_name}: {e}")
        message_history = [message] # Fallback
        logger.warning("Warning: Proceeding with AI generation using only the triggering message.")

    # --- Generate AI Response ---
    ai_response = await generate_ai_response(channel_name, message_history, message)

    if not ai_response:
        logger.warning("AI generation failed or returned empty.")
        return

    # --- Send Response ---
    try:
        # Slightly shorter delay might feel okay given it *always* responds
        delay = random.uniform(0.8, 2.5)
        await asyncio.sleep(delay)

        # Typing indicator can still be nice
        async with message.channel.typing():
            typing_delay = random.uniform(0.5, 1.5)
            await asyncio.sleep(typing_delay)

        await message.channel.send(ai_response)
        logger.info(f"Responded in #{channel_name} to {message.author.name}: '{ai_response}'")

    except discord.errors.Forbidden: # Should be caught by check above, but good safety
        logger.error(f"Permission error sending message in #{channel_name} despite initial check.")
    except Exception as e:
        logger.exception(f"An unexpected error occurred during response sending in #{channel_name}: {e}")

# --- Run the Bot ---
if __name__ == "__main__":
    # Check essential config before trying to run
    if not BOT_TOKEN:
        logger.critical("CRITICAL: BOT_TOKEN not found in environment variables. Bot cannot start.")
    # Note: We check azure_client availability within the event handlers now
    else:
        try:
            logger.info("Starting Always Responder Discord bot...")
            discord_client.run(BOT_TOKEN)
        except discord.errors.LoginFailure:
            logger.critical("CRITICAL: Failed to log in to Discord. Check BOT_TOKEN.")
        except Exception as e:
            logger.critical(f"CRITICAL: An error occurred while running the Discord bot: {e}", exc_info=True)