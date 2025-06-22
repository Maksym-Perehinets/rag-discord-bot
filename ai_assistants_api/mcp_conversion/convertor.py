import os
import traceback
from contextlib import AsyncExitStack
from typing import Optional

from mcp import ClientSession
from mcp.client.streamable_http import streamablehttp_client
from mcp.types import CallToolResult

from logger import logger


class MCPConvertor:
    def __init__(self):
        self.session: Optional[ClientSession] = None
        self.exit_stack = AsyncExitStack()
        self.mcp_server_url = os.getenv("MCP_SERVER_URL", "http://localhost:8080/mcp")

    async def connect_to_server(self):
        logger.info(f"Setting up connection to mcp server on following host {self.mcp_server_url}")
        try:
            mcp_session_manager = streamablehttp_client(url=self.mcp_server_url)
            read_stream, write_stream, get_session_id = await self.exit_stack.enter_async_context(mcp_session_manager)

            session_manager = ClientSession(read_stream, write_stream)
            self.session = await self.exit_stack.enter_async_context(session_manager)

            await self.session.initialize()

            # List available tools
            response = await self.session.list_tools()
            tools = response.tools
            logger.info("\nConnected to server with tools: %s", [tool.name for tool in tools])
        except Exception as e:
            logger.fatal(f"Failed to connect to MCP server: {e}")
            raise

    async def get_converted_tools(self) -> list[dict[str, str | dict[str, any]]]:
        logger.info(f"Getting and converting available tools in mcp server on following host {self.mcp_server_url}")
        try:
            response = await self.session.list_tools()
        except Exception as e:
            logger.fatal(f"Cannot get available toolf from following mcp server {self.mcp_server_url} exited with following error: {str(e)}")

        # Create a serializable representation of tools for the OpenAI API
        available_tools = [{
            "type": "function",
            "function": {
                "name": tool.name,
                "description": tool.description,
                "parameters": tool.inputSchema
            }
        } for tool in response.tools]

        return available_tools

    async def execute_tool(self, tool_name: str, tool_args: dict) -> CallToolResult:
        try:
            return await self.session.call_tool(tool_name, tool_args)
        except Exception as e:
            logger.error(f"Failed to execute tool {tool_name} with args {tool_args}: {e}")
            logger.error(traceback.format_exc())

    async def cleanup(self):
        """Clean up resources"""
        logger.info("\nClosing connections...")
        await self.exit_stack.aclose()
