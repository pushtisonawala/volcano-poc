# NamespaceQueue PoC for Volcano Scheduler

## Overview
This prototype demonstrates a namespace-scoped `NamespaceQueue` CustomResourceDefinition (CRD) for the Volcano batch scheduler. Unlike the default cluster-scoped `Queue`, a `NamespaceQueue` allows tenants to create and manage their own queues within their namespace, referencing a parent cluster `Queue` for resource limits.

## Problem Solved
- Enables multi-tenancy: Namespace users can create their own queues without cluster-admin permissions.
- Resource isolation: Each `NamespaceQueue` is limited by its parent cluster `Queue`.
- No impact on existing cluster `Queue` behavior.

## Relationship to Cluster Queue
- `NamespaceQueue` is namespace-scoped and always references a parent cluster-scoped `Queue` via the `parentQueue` field.
- Resources in a `NamespaceQueue` cannot exceed the parent `Queue`'s capability.
- PodGroups can reference a `NamespaceQueue` using the same annotation: `scheduling.volcano.sh/queue-name: <namespacequeue-name>`

## How to Run Locally (with kind)
1. Install [kind](https://kind.sigs.k8s.io/) and [kubectl](https://kubernetes.io/docs/tasks/tools/).
2. Deploy Volcano scheduler and CRDs.
3. Apply the `NamespaceQueue` CRD:
   ```sh
   kubectl apply -f config/crd/namespacequeue.yaml
   ```
4. Deploy the controller:
   ```sh
   go run main.go
   ```
5. Create example queues:
   ```sh
   kubectl apply -f examples/team-a-queue.yaml
   kubectl apply -f examples/team-b-queue.yaml
   ```

## Design Decisions
- `NamespaceQueue` is namespace-scoped for tenant self-service.
- Always references a parent cluster `Queue`.
- Controller validates parent existence and updates status accordingly.
- Finalizer ensures cleanup.
- No changes to existing cluster `Queue` logic.

## Full Implementation Would Add
- Webhook validation for resource limits.
- Integration with Volcano scheduling logic.
- RBAC for fine-grained access.
- More status fields and events.
- E2E tests and CI/CD.
