package vogcluster

// HeaderClaimerHash is the NATS message header name used by game
// instances and the coordinator to defend against two processes
// claiming the same instance_id. The value is an opaque hash derived
// from the instance's pod UID (or equivalent immutable identity). The
// coordinator records the first hash seen for a given instance_id and
// rejects any subsequent register/heartbeat/deregister with a different
// hash.
const HeaderClaimerHash = "X-Claimer-Hash"
