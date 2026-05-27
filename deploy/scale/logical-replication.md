# Postgres Logical Replication Recipe

Use logical replication when:

- Promoting a region from "warm secondary" to "active replica".
- Keeping the research-API sandbox replica fresh from the primary.

## Steps

1. On the primary, set `wal_level=logical`, `max_replication_slots=20`, and `max_wal_senders=20`.
2. Create a publication on the primary:
   ```sql
   CREATE PUBLICATION joblantern_pub FOR TABLE verifications, rental_verifications, edu_verifications;
   ```
3. On the secondary, create the subscription:
   ```sql
   CREATE SUBSCRIPTION joblantern_sub
     CONNECTION 'host=primary dbname=joblantern user=replicator password=...'
     PUBLICATION joblantern_pub;
   ```
4. Verify with `pg_stat_subscription` that `last_msg_send_time` ticks every second.

## Gotchas

- Sequences are **not** replicated logically — use `ALTER SEQUENCE ... RESTART` after failover.
- DDL is not replicated. Schema changes ship through the regular goose migration pipeline.
- The research-sandbox replica strips the `submitter_id` column at the application layer before exposing rows.
