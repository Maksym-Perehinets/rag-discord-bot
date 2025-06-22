from fastapi import APIRouter, Depends

from llms.azure_openai import LLMClient
from ..models.common import CommonRequest, CommonResponse
from ..services import azure_openai
from ..dependencies import get_llm_client

router = APIRouter(
    prefix="/azure-openai",
)

@router.post("/process-query")
async def process_query(request: CommonRequest, llm_client: LLMClient = Depends(get_llm_client)) -> CommonResponse:
    """
    Endpoint to process a query using Azure OpenAI.
    """

    response = await azure_openai.request_to_llm(llm_client, request.query)

    return CommonResponse(
        user_id=request.user_id,
        message_id=request.message_id,
        query_response=response
    )