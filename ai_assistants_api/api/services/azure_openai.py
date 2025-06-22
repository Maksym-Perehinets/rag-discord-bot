from llms.azure_openai import LLMClient
from logger import logger


async def request_to_llm(ai_client: LLMClient, query: list[dict[str, str]]) -> dict[str, str]:
    """
    Function to process a request using the Azure OpenAI LLM client.
    This function simulates the processing of a query and returns a response.
    """
    logger.info(f"Processing query: {query}")

    try:
        response = await ai_client.process_query(query)
        logger.info(f"LLM response: {response}")
        return {"response": response}
    except Exception as e:
        logger.error(f"Error processing query with LLM: {e}")
        return {"error": "Error processing query with LLM."}
