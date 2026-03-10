---
layout: home

hero:
  name: DriftGuard
  text: API Type Safety
  tagline: Enforce API type safety across OpenAPI, GraphQL, and gRPC — catch breaking changes before they reach production.
  actions:
    - theme: brand
      text: Get Started
      link: /install
    - theme: alt
      text: npm SDK
      link: /npm
    - theme: alt
      text: GitHub Marketplace
      link: https://github.com/marketplace/actions/drift-guard
    - theme: alt
      text: View on GitHub
      link: https://github.com/pgomes13/drift-guard-engine

features:
  - title: Multi-schema support
    details: Parses OpenAPI 3.x (YAML/JSON), GraphQL SDL, and Protobuf (.proto) schemas.
  - title: Severity classification
    details: Every change is classified as breaking, non-breaking, or info — with detailed rules per schema type.
  - title: CI-ready
    details: Posts PR comments with the full diff, updates a drift log on GitHub Pages, and supports --fail-on-breaking to block merges.
  - title: SDK & npm package
    details: Use drift-guard programmatically from Node.js/TypeScript via @pgomes13/drift-guard, or import pkg/compare and pkg/impact directly in Go.
---
