// Package vogcluster defines the message types and NATS subject topology
// used to coordinate the horizontally-scaled vog game cluster.
//
// This package contains only data types, subject builders, and encoding
// helpers. It does not depend on nats.go in its public API — consumers
// pass *nats.Conn from their own services.
//
// See docs/superpowers/specs/2026-04-13-horizontal-scaling-design.md in
// the vog-arch repository for the full architecture.
package vogcluster
