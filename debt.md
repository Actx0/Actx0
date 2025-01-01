# File review checklist

Manually validate each file for AI slop, unnecessary defensive code, and style drift.

Mark done: `- [x] path/to/file`

---

## Root

- [x] actx0.go

## api/

- [ ] api/access.go
- [ ] api/agent.go
- [ ] api/api_key.go
- [ ] api/billing.go
- [ ] api/customer.go
- [x] api/health.go
- [x] api/health_test.go
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
- [x] api/ready.go
- [ ] api/register.go
- [ ] api/session.go
- [x] api/settings.go
- [x] api/setup.go
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
- [x] core/server.go

## db/

- [ ] db/agent.go
- [ ] db/atask.go
- [x] db/audit.go
- [x] db/config.go
- [x] db/database.go
- [x] db/helpers.go
- [x] db/integration.go
- [x] db/invite.go
- [ ] db/memory.go
- [ ] db/message.go
- [x] db/preset.go
- [ ] db/prompt.go
- [ ] db/session.go
- [x] db/subscription.go
- [x] db/types.go
- [x] db/usage.go
- [x] db/user.go
- [x] db/user_api_key.go
- [x] db/user_session.go
- [ ] db/wdoc.go
- [x] db/wkey.go
- [x] db/workspace.go
- [x] db/wstats.go
- [x] db/wuser.go

## locale/

- [x] locale/locale.go

## mcp/

- [x] mcp/memory.go
- [x] mcp/server.go

## middleware/

- [x] middleware/auth.go
- [x] middleware/ctx.go
- [x] middleware/helpers.go
- [x] middleware/logger.go
- [x] middleware/prometheus.go
- [x] middleware/protect.go
- [x] middleware/request_id.go
- [x] middleware/request_size.go
- [x] middleware/usage.go

## migration/

- [x] migration/migration.go
- [x] migration/qdb.go
- [x] migration/registry.go

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

- [x] pkg/ai/chat.go
- [x] pkg/ai/embed.go
- [x] pkg/ai/errors.go
- [x] pkg/ai/usage.go

## pkg/async/

- [x] pkg/async/pool.go

## pkg/cache/

- [x] pkg/cache/cache.go
- [x] pkg/cache/config.go
- [x] pkg/cache/errors.go
- [x] pkg/cache/redis.go

## pkg/chunk/

- [x] pkg/chunk/chunk.go

## pkg/event/

- [x] pkg/event/event.go

## pkg/meteo/

- [x] pkg/meteo/client.go
- [x] pkg/meteo/errors.go
- [x] pkg/meteo/types.go

## pkg/qdrant/

- [x] pkg/qdrant/client.go
- [x] pkg/qdrant/errors.go
- [x] pkg/qdrant/helpers.go
- [x] pkg/qdrant/types.go

## pkg/resend/

- [x] pkg/resend/client.go
- [x] pkg/resend/helpers.go

## pkg/storage/

- [x] pkg/storage/config.go
- [x] pkg/storage/document.go
- [x] pkg/storage/errors.go
- [x] pkg/storage/local.go
- [x] pkg/storage/s3.go
- [x] pkg/storage/store.go

## pkg/stripe/

- [x] pkg/stripe/checkout.go
- [x] pkg/stripe/client.go
- [x] pkg/stripe/config.go
- [x] pkg/stripe/errors.go
- [x] pkg/stripe/portal.go
- [x] pkg/stripe/subscription.go
- [x] pkg/stripe/webhook.go

## pkg/util/

- [x] pkg/util/cookies.go
- [x] pkg/util/cookies_test.go
- [x] pkg/util/file_sys.go
- [x] pkg/util/file_sys_test.go
- [x] pkg/util/fvalidate.go
- [x] pkg/util/gavatar.go
- [x] pkg/util/gavatar_test.go
- [x] pkg/util/handle_test.go
- [x] pkg/util/hashing.go
- [x] pkg/util/hashing_test.go
- [x] pkg/util/helpers.go
- [x] pkg/util/loader.go
- [x] pkg/util/loader_test.go
- [x] pkg/util/validation.go
- [x] pkg/util/validation_test.go
- [x] pkg/util/validators.go

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

## web/

- [ ] web/index.html
- [ ] web/postcss.config.js
- [ ] web/tailwind.config.js
- [ ] web/vite.config.js
- [ ] web/.eslintrc.cjs

## web/src/

- [ ] web/src/App.vue
- [ ] web/src/i18n.js
- [ ] web/src/main.js
- [ ] web/src/style.css

## web/src/api/

- [ ] web/src/api/index.js
- [ ] web/src/api/session-resources.js

## web/src/components/

- [ ] web/src/components/AppNav.vue
- [ ] web/src/components/CaseLayout.vue
- [ ] web/src/components/CookieConsent.vue
- [ ] web/src/components/EmptyPage.vue
- [ ] web/src/components/FlashNotice.vue
- [ ] web/src/components/LandingNav.vue
- [ ] web/src/components/PromptFormFields.vue

## web/src/lib/

- [ ] web/src/lib/agent.js
- [ ] web/src/lib/auth.js
- [ ] web/src/lib/code.js
- [ ] web/src/lib/cookies.js
- [ ] web/src/lib/flash.js
- [ ] web/src/lib/landingPage.js
- [ ] web/src/lib/pagination.js
- [ ] web/src/lib/permission.js
- [ ] web/src/lib/preferences.js
- [ ] web/src/lib/prompt.js

## web/src/mocks/

- [ ] web/src/mocks/agents.js
- [ ] web/src/mocks/integrations.js

## web/src/router/

- [ ] web/src/router/index.js

## web/src/utils/

- [ ] web/src/utils/date.js
- [ ] web/src/utils/storage.js

## web/src/views/

- [ ] web/src/views/Agent.vue
- [ ] web/src/views/Agents.vue
- [ ] web/src/views/Audits.vue
- [ ] web/src/views/Billing.vue
- [ ] web/src/views/Dashboard.vue
- [ ] web/src/views/Docs.vue
- [ ] web/src/views/ForgotPassword.vue
- [ ] web/src/views/Integrations.vue
- [ ] web/src/views/Invite.vue
- [ ] web/src/views/Knowledge.vue
- [ ] web/src/views/Landing.vue
- [ ] web/src/views/Login.vue
- [ ] web/src/views/Members.vue
- [ ] web/src/views/NotFound.vue
- [ ] web/src/views/Profile.vue
- [ ] web/src/views/PromptDetail.vue
- [ ] web/src/views/Prompts.vue
- [ ] web/src/views/Register.vue
- [ ] web/src/views/ResetPassword.vue
- [ ] web/src/views/SelectWorkspace.vue
- [ ] web/src/views/ServerError.vue
- [ ] web/src/views/Session.vue
- [ ] web/src/views/Settings.vue
- [ ] web/src/views/Setup.vue
- [ ] web/src/views/Status.vue
- [ ] web/src/views/VerifyEmail.vue
- [ ] web/src/views/Workspaces.vue

## web/src/views/cases/

- [ ] web/src/views/cases/CustomerSupport.vue
- [ ] web/src/views/cases/DevTools.vue
- [ ] web/src/views/cases/ECommerce.vue
- [ ] web/src/views/cases/Education.vue
- [ ] web/src/views/cases/Healthcare.vue
- [ ] web/src/views/cases/Sales.vue
