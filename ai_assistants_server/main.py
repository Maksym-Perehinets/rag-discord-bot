from contextlib import asynccontextmanager

from fastapi import FastAPI

from api import startup
from llms import azure_openai
from logger import logger
from mcp_conversion import convertor


@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Lifespan context manager for the FastAPI application.
    This is used to initialize resources when the app starts and clean them up when it stops.
    """
    logger.info("Starting up the application...")

    mcp = convertor.MCPConvertor()
    try:
        await mcp.connect_to_server()
        logger.info("Connected to MCP server successfully.")
        tools = await mcp.get_converted_tools()
        logger.info(f"Available tools: {tools}")
    except Exception as e:
        logger.fatal(f"Failed to connect to MCP server: {e}")
        raise

    tool_executor = lambda tool_name, tool_args: mcp.execute_tool(tool_name, tool_args)

    llm_client = azure_openai.LLMClient(tools, tool_executor)


    app.state.mcp = mcp
    app.state.llm_client = llm_client

    logger.info("Startup complete. Application is ready to serve requests.")

    yield

    await mcp.cleanup()


app = startup.get_app(lifespan)


