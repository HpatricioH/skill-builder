import openai
import os
from typing import Optional

def call_gpt(prompt: str, model: str = "gpt-3.5-turbo", max_tokens: int = 150, temperature: float = 0.7) -> Optional[str]:
    """
    Call GPT API to generate a response based on the given prompt.
    
    Args:
        prompt (str): The input prompt for GPT
        model (str): The GPT model to use (default: gpt-3.5-turbo)
        max_tokens (int): Maximum number of tokens in the response
        temperature (float): Controls randomness (0.0 to 1.0)
    
    Returns:
        str: The GPT response, or None if there was an error
    """
    try:
        # You'll need to set your OpenAI API key
        # You can do this by setting the OPENAI_API_KEY environment variable
        # or by uncommenting and setting the line below:
        # openai.api_key = "your-api-key-here"
        
        client = openai.OpenAI(
            api_key=os.getenv("OPENAI_API_KEY")
        )
        
        response = client.chat.completions.create(
            model=model,
            messages=[
                {"role": "user", "content": prompt}
            ],
            max_tokens=max_tokens,
            temperature=temperature
        )
        
        return response.choices[0].message.content.strip()
        
    except Exception as e:
        print(f"Error calling GPT: {e}")
        return None

def generate_recipe(ingredients: list) -> Optional[str]:
    """
    Generate a recipe based on available ingredients.
    
    Args:
        ingredients (list): List of available ingredients
    
    Returns:
        str: A recipe suggestion, or None if there was an error
    """
    ingredients_str = ", ".join(ingredients)
    prompt = f"Create a simple recipe using these ingredients: {ingredients_str}. Include preparation steps."
    
    return call_gpt(prompt, max_tokens=300)

def suggest_groceries(current_items: list) -> Optional[str]:
    """
    Suggest additional groceries based on current pantry items.
    
    Args:
        current_items (list): List of current pantry items
    
    Returns:
        str: Grocery suggestions, or None if there was an error
    """
    items_str = ", ".join(current_items)
    prompt = f"Based on these pantry items: {items_str}, suggest 5 additional grocery items that would complement them well for cooking."
    
    return call_gpt(prompt, max_tokens=200)

