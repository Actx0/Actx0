# Go file review checklist

Manually validate each file for AI slop, unnecessary defensive code, and style drift.

Mark done: `- [x] path/to/file.go`

**Progress:** 14 / 171

---

## Root

- [x] actx0.go

## api/

- [ ] api/access.go
- [ ] api/agent.go
- [ ] api/api_key.go
- [ ] api/billing.go
- [ ] api/customer.go
- [ ] api/health.go
- [ ] api/health_test.go
- [ ] api/integration.go
- [ ] api/invite.go
- [ ] api/knowledge.go
- [ ] api/login.go
- [x] api/logout.go
- [ ] api/me.go
- [ ] api/member.go
- [ ] api/memory.go
- [ ] api/message.go
- [ ] api/oauth.go
- [ ] api/password.go
- [ ] api/profile.go
- [ ] api/prompt.go
- [ ] api/ready.go
- [ ] api/register.go
- [ ] api/session.go
- [ ] api/settings.go
- [ ] api/setup.go
- [ ] api/stats.go
- [ ] api/subscription.go
- [ ] api/tools.go
- [ ] api/workspace.go

## cli/

- [x] cli/license.go
- [x] cli/mcp.go
- [x] cli/migrate.go
- [x] cli/root.go
- [x] cli/server.go
- [x] cli/version.go

## core/

- [x] core/database.go
- [x] core/loader.go
- [x] core/logger.go
- [ ] core/server.go

## db/

- [ ] db/agent.go
- [ ] db/atask.go
- [x] db/audit.go
- [x] db/config.go
- [x] db/database.go
- [ ] db/helpers.go
- [ ] db/integration.go
- [ ] db/invite.go
- [ ] db/memory.go
- [ ] db/message.go
- [ ] db/preset.go
- [ ] db/prompt.go
- [ ] db/session.go
- [ ] db/subscription.go
- [ ] db/types.go
- [ ] db/usage.go
- [ ] db/user.go
- [ ] db/user_api_key.go
- [ ] db/user_session.go
- [ ] db/wdoc.go
- [ ] db/wkey.go
- [ ] db/workspace.go
- [ ] db/wstats.go
- [ ] db/wuser.go

## locale/

- [ ] locale/locale.go

## mcp/

- [ ] mcp/memory.go
- [ ] mcp/server.go

## middleware/

- [ ] middleware/auth.go
- [ ] middleware/basicauth.go
- [ ] middleware/ctx.go
- [ ] middleware/helpers.go
- [ ] middleware/logger.go
- [ ] middleware/permission.go
- [ ] middleware/prometheus.go
- [ ] middleware/request_id.go
- [ ] middleware/request_size.go
- [ ] middleware/role.go
- [ ] middleware/usage.go
- [ ] middleware/workspace.go

## migration/

- [ ] migration/migration.go
- [ ] migration/qdb.go
- [ ] migration/registry.go

## module/

- [ ] module/access.go
- [ ] module/agent.go
- [ ] module/api_key.go
- [ ] module/asession.go
- [ ] module/async.go
- [ ] module/auth.go
- [ ] module/billing.go
- [ ] module/ctx.go
- [ ] module/document.go
- [ ] module/events.go
- [ ] module/invite.go
- [ ] module/me.go
- [ ] module/memory.go
- [ ] module/message.go
- [ ] module/permission.go
- [ ] module/profile.go
- [ ] module/prompt.go
- [ ] module/register.go
- [ ] module/session.go
- [ ] module/settings.go
- [ ] module/setup.go
- [ ] module/stats.go
- [ ] module/tools.go
- [ ] module/usage.go
- [ ] module/wspace.go
- [ ] module/wuser.go

## pkg/ai/

- [ ] pkg/ai/chat.go
- [ ] pkg/ai/embed.go
- [ ] pkg/ai/errors.go
- [ ] pkg/ai/usage.go

## pkg/async/

- [ ] pkg/async/pool.go

## pkg/cache/

- [ ] pkg/cache/cache.go
- [ ] pkg/cache/config.go
- [ ] pkg/cache/errors.go
- [ ] pkg/cache/redis.go

## pkg/chunk/

- [ ] pkg/chunk/chunk.go

## pkg/event/

- [ ] pkg/event/event.go

## pkg/meteo/

- [ ] pkg/meteo/client.go
- [ ] pkg/meteo/errors.go
- [ ] pkg/meteo/types.go

## pkg/qdrant/

- [ ] pkg/qdrant/client.go
- [ ] pkg/qdrant/errors.go
- [ ] pkg/qdrant/helpers.go
- [ ] pkg/qdrant/types.go

## pkg/resend/

- [ ] pkg/resend/client.go
- [ ] pkg/resend/helpers.go

## pkg/storage/

- [ ] pkg/storage/config.go
- [ ] pkg/storage/document.go
- [ ] pkg/storage/errors.go
- [ ] pkg/storage/local.go
- [ ] pkg/storage/s3.go
- [ ] pkg/storage/store.go

## pkg/stripe/

- [ ] pkg/stripe/checkout.go
- [ ] pkg/stripe/client.go
- [ ] pkg/stripe/config.go
- [ ] pkg/stripe/errors.go
- [ ] pkg/stripe/portal.go
- [ ] pkg/stripe/subscription.go
- [ ] pkg/stripe/webhook.go

## pkg/util/

- [ ] pkg/util/cookies.go
- [ ] pkg/util/cookies_test.go
- [ ] pkg/util/file_sys.go
- [ ] pkg/util/file_sys_test.go
- [ ] pkg/util/fvalidate.go
- [ ] pkg/util/gavatar.go
- [ ] pkg/util/gavatar_test.go
- [ ] pkg/util/handle_test.go
- [ ] pkg/util/hashing.go
- [ ] pkg/util/hashing_test.go
- [ ] pkg/util/helpers.go
- [ ] pkg/util/loader.go
- [ ] pkg/util/loader_test.go
- [ ] pkg/util/validation.go
- [ ] pkg/util/validation_test.go
- [ ] pkg/util/validators.go

## service/agent/

- [ ] service/agent/context.go
- [ ] service/agent/errors.go
- [ ] service/agent/memory.go
- [ ] service/agent/service.go
- [ ] service/agent/types.go

## service/knowledge/

- [ ] service/knowledge/delete.go
- [ ] service/knowledge/helpers.go
- [ ] service/knowledge/index.go
- [ ] service/knowledge/search.go
- [ ] service/knowledge/service.go

## task/

- [ ] task/agent.go
- [ ] task/doc.go
- [ ] task/task.go
- [ ] task/util.go
