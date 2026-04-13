Help create or review a database migration for the Go/GORM project.

## For creating a new migration

1. **Understand the change**: What models are being added/modified? Parse from arguments.
2. **Generate the model changes**:
   - Add/modify GORM model structs in `internal/{domain}/model/`
   - Use proper GORM tags: `gorm:"column:name;type:varchar(255);not null;uniqueIndex"`
   - Pointer types for nullable fields (`*string`, `*time.Time`)
   - Add `BeforeCreate` hook for UUID generation if new model
3. **Update AutoMigrate**: Add new models to `cmd/api/main.go` AutoMigrate call.
4. **Handle relationships**:
   - Foreign keys: `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
   - Preload configuration in repository layer
5. **Check data integrity**:
   - Will this break existing data? Are there default values needed?
   - Does the migration need to backfill data?
   - Are there indexes needed for query performance?
6. **Update downstream**:
   - Repository methods for new fields/associations
   - DTOs for request/response changes
   - Service layer business logic
   - Regenerate mocks if interfaces changed

## For reviewing a migration

- Check for destructive operations (column drops, type changes)
- Verify indexes exist for foreign keys and frequently queried columns
- Check constraint naming for clarity
- Verify nullable vs non-nullable is intentional
- Test with `go build ./...` to verify model compiles

$ARGUMENTS
