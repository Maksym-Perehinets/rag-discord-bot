from pydantic import BaseModel

class CommonResponse(BaseModel):
    """
    Common response model for API responses.
    Contains the user ID, message ID, and a dictionary of query responses.
    """
    user_id: str
    message_id: str
    query_response: dict[str, str]

class CommonRequest(BaseModel):
    """
    Common request model for API requests.
    Contains the user ID, message ID, and a dictionary of query parameters.
    """
    user_id: str
    message_id: str
    query: list[dict[str, str]]