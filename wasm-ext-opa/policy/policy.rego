package envoy.authz

import data.secret.secret
import input

default allow = false

token = {"valid": valid, "payload": payload} {
    [_, encoded] := split(input.authorization, " ")
    [valid, _, payload] := io.jwt.decode_verify(encoded, {"secret": secret})
}

allow {
  is_token_valid
  # action_allowed
}

is_token_valid {
    token.valid
    now := time.now_ns() / 1000000000
    token.payload.nbf <= now
    now < token.payload.exp
}

