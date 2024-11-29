import re
import json

def clean_string(text):
    
    """
    This function takes in a string and performs a series of text cleaning operations.

    Args:
        text (str): The text to be cleaned. This is expected to be a string.

    Returns:
        cleaned_text (str): The cleaned text after all the cleaning operations
        have been performed.
    """
    # Stripping and reducing multiple spaces to single:
    cleaned_text = re.sub(r"\s+", " ", text.strip())

    # Removing backslashes:
    cleaned_text = cleaned_text.replace("\\", "")

    # Replacing hash characters:
    cleaned_text = cleaned_text.replace("#", " ")

    # Eliminating consecutive non-alphanumeric characters:
    # This regex identifies consecutive non-alphanumeric characters (i.e., not
    # a word character [a-zA-Z0-9_] and not a whitespace) in the string
    # and replaces each group of such characters with a single occurrence of
    # that character.
    # For example, "!!! hello !!!" would become "! hello !".
    cleaned_text = re.sub(r"([^\w\s])\1*", r"\1", cleaned_text)

    return cleaned_text


# check if the source is valid json string
def is_valid_json_string(source: str):
    try:
        _ = json.loads(source)
        return True
    except json.JSONDecodeError:
        return False
    

def tw2s(text):
    '''
        台湾繁体转大陆简体：https://pypi.org/project/OpenCC/
    '''
    from opencc import OpenCC
    converter = OpenCC('tw2s.json')
    text = converter.convert(text)
    return text
    