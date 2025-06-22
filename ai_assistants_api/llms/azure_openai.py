from openai import AsyncAzureOpenAI
from logger import logger

import asyncio
import json
import os


class LLMClient:
    def __init__(self, tools: list[dict[str, str | dict[str, any]]], tool_executor: callable([str, dict]), model: str = "gpt-4o"):
        logger.info(f"Initializing LLMClient with model: {model} and tools: {tools}")
        self.model = model
        self.tools = tools
        self.tools_executor = tool_executor

        self.client = AsyncAzureOpenAI(
            azure_endpoint=os.getenv("AZURE_OPENAI_ENDPOINT"),
            api_key=os.getenv("AZURE_OPENAI_API_KEY"),
            api_version="2024-05-01-preview"
        )

    async def process_query(self, query: list[dict[str, str]]) -> str:
        messages = query
        max_iterations = 10

        for i in range(max_iterations):
            logger.info(f"--- Starting LLM Iteration {i + 1} ---")

            response = await self.client.chat.completions.create(
                model=self.model,
                messages=messages,
                tools=self.tools,
                tool_choice="auto",
            )

            response_message = response.choices[0].message
            logger.info(f"Received response from LLM: {response_message}")

            messages.append(response_message.model_dump(exclude_unset=True))

            if not response_message.tool_calls:
                logger.info("LLM provided a final answer. Ending iterations.")
                return response_message.content or "The agent did not provide a final text response."

            logger.info("LLM requested tool calls. Executing now.")
            for tool_call in response_message.tool_calls:
                tool_name = tool_call.function.name

                try:
                    tool_args = json.loads(tool_call.function.arguments) if tool_call.function.arguments else {}
                except json.JSONDecodeError:
                    logger.warning(f"Warning: Failed to parse JSON arguments for tool {tool_name}. "
                                   f"Arguments: {tool_call.function.arguments}")
                    tool_args = {}

                result = await self.tools_executor(tool_name, tool_args)

                tool_content = json.dumps(result, default=str)

                logger.info(f"Tool `{tool_name}` result: {tool_content}")

                messages.append({
                    "tool_call_id": tool_call.id,
                    "role": "tool",
                    "name": tool_name,
                    "content": tool_content,
                })

        logger.warning(f"Reached max iterations ({max_iterations}) without a final text response.")
        return "The agent could not produce a final answer within the allowed number of iterations."

    async def chat_loop(self):
        """Local test only for now=)"""
        print("\nMCP Client Started!")
        print("Type your queries or 'quit' to exit.")

        while True:
            try:
                query = await asyncio.to_thread(input, "\nQuery: ")
                query = query.strip()

                if query.lower() == 'quit':
                    break

                if not query:
                    continue

                response = await self.process_query(query)
                print("\n" + response)

            except (KeyboardInterrupt, EOFError):
                print("\nExiting...")
                break
            except Exception as e:
                print(f"\nAn unexpected error occurred: {e}")
