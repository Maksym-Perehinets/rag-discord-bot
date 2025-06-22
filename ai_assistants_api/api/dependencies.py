from fastapi import Request, HTTPException, status
from mcp_conversion.convertor import MCPConvertor
from llms.azure_openai import LLMClient

# TODO remove is not used and will most likely be obsolete in the future
# def get_mcp_convertor(request: Request) -> MCPConvertor:
#     """
#     Dependency provider for the MCPConvertor client.
#     Gets the shared instance from the application state.
#     """
#     mcp_client = request.app.state.mcp_convertor
#     if not mcp_client:
#         raise HTTPException(
#             status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
#             detail="MCP client is not available or connected."
#         )
#     return mcp_client

def get_llm_client(request: Request) -> LLMClient:
    """
    Dependency provider for the LLMClient.
    Gets the shared instance from the application state.
    """
    llm_client = request.app.state.llm_client
    if not llm_client:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="LLM client is not available."
        )
    return llm_client