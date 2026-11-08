// Package web serves the lab: server-rendered HTML, hand-written CSS, a
// little vanilla JS, and live vitals over SSE.
//
// There is no frontend framework, no build step, and no npm anywhere in this
// project. If a page needs JavaScript it is a hand-written script served by
// the lab, and the app works offline.
package web

// Routes served by the lab's human-facing surface.
const (
	RouteHome     = "/"              // the fleet: one creature card per familiar
	RouteFamiliar = "/familiar/{id}" // the lab file for one familiar
	RouteEvents   = "/events"        // SSE stream of live vitals
	RouteStatic   = "/static/"       // hand-written CSS and JS
)

// UnreachableAfter is how long a familiar may go without reporting before its
// card shows an unreachable / sleeping state.
const UnreachableAfter = 60 // seconds

// TODO(phase-1, task 6.1): the home page — server-rendered creature cards
// with inline SVG and vitals, dark storm-lab theme.
// TODO(phase-1, task 6.2): the SSE hub and GET /events; vanilla JS patches
// vitals in place; cards go unreachable after UnreachableAfter seconds.
// TODO(phase-1, task 6.3): the familiar detail page stub.
