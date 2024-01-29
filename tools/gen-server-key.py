from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.serialization import (
    Encoding,
    PublicFormat,
)
import base64

private_key = ec.generate_private_key(ec.SECP256R1)
private_key_numbers = private_key.private_numbers().private_value

public_key = private_key.public_key()
public_key_bytes = public_key.public_bytes(Encoding.X962, PublicFormat.UncompressedPoint)
public_key_decoded = base64.urlsafe_b64encode(public_key_bytes).decode()

print('private key: %s' % (private_key_numbers))
print(' public key: %s' % (public_key_decoded))



