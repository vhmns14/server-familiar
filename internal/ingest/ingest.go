// Package ingest holds the lab's agent-facing HTTP surface: enrolment and
// metric ingest.
//
// Agents are outbound-only. The lab never connects into a server; it answers
// enrolment, accepts samples, and hands out work for agents to poll.
package ingest

// Header and route names shared by the agent and the lab.
const (
	// TokenHeader carries the fleet token during enrolment and the
	// per-familiar token on every sample afterwards.
	TokenHeader = "Authorization"

	RouteEnrol  = "/api/enrol"
	RouteIngest = "/api/ingest"
)

// TODO(phase-1, task 4.1): POST /api/enrol — fleet token, derives species from
// the host profile, creates familiar + first life, returns a per-familiar
// token (random 32 bytes, stored hashed).
// TODO(phase-1, task 4.2): POST /api/ingest — per-familiar token, stores the
// sample, touches last-seen, publishes to the SSE hub, 401 on a bad token.
