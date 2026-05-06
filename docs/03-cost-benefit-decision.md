# Architecture Decision Record

## ADR-001: WhatsApp Provider Selection

## Status

Accepted

## Context

The system requires:

* production reliability
* predictable maintenance
* reasonable operational costs

## Alternatives Considered

1. WhatsApp Cloud API
2. WhatsMeow
3. Twilio
4. Zenvia
5. 360Dialog

## Decision

Use WhatsApp Cloud API.

## Decision Drivers

* official provider
* webhook reliability
* lower maintenance burden
* scalability
* long-term viability

## Tradeoffs

### Positive

* robust integration
* safer production environment
* future compatibility

### Negative

* moderate cost
* onboarding setup

## Final Conclusion

Cloud API provides the best long-term ROI.

