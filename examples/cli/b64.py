import base64
import sys

method = sys.argv[1]
strinput = sys.argv[2]

if method == "encode":
    string_bytes = strinput.encode("ascii")
    
    base64_bytes = base64.b64encode(string_bytes)
    result = base64_bytes.decode("ascii")
    print(result)

if method == "decode":
    base64_bytes = strinput.encode("ascii")
  
    string_bytes = base64.b64decode(base64_bytes)
    result = string_bytes.decode("ascii")
    print(result)