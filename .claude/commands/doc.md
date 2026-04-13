Generate or update Swagger/OpenAPI annotations for Go HTTP handlers.

## Instructions

1. **Read the handler**: Understand the endpoint — method, path, request body, query params, path params, response.

2. **Generate Swagger comments** above the handler function:
   ```go
   // @Summary      Short description of what this endpoint does
   // @Description  Longer description if needed
   // @Tags         domain-name
   // @Accept       json
   // @Produce      json
   // @Param        id    path     string  true  "Resource ID"
   // @Param        body  body     dto.CreateRequest  true  "Request body"
   // @Param        page  query    int     false "Page number" default(1)
   // @Param        limit query    int     false "Page size"   default(20)
   // @Success      200   {object} dto.Response
   // @Success      201   {object} dto.Response
   // @Failure      400   {object} map[string]string
   // @Failure      401   {object} map[string]string
   // @Failure      404   {object} map[string]string
   // @Failure      500   {object} map[string]string
   // @Security     BearerAuth
   // @Router       /api/v1/resources/{id} [get]
   ```

3. **Rules**:
   - `@Tags` must match the domain name (product, user, order, cart)
   - `@Security BearerAuth` only on authenticated endpoints
   - Include all possible `@Failure` codes the handler can return
   - `@Param` for every path param, query param, and request body
   - Use actual DTO struct types, not `interface{}`

4. **Regenerate docs**: Run `make doc` (`swag fmt && swag init -g ./cmd/api/main.go`)

5. **Verify**: Check that `docs/swagger.json` was updated correctly.

$ARGUMENTS
