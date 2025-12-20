#!/bin/bash

# Script untuk menjalankan migration perubahan foreign key constraint
# File: apply_fk_fix.sh

echo "🔧 Applying Foreign Key Constraint Fix..."
echo "=========================================="
echo ""

# Database connection details
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-your_database_name}"
DB_USER="${DB_USER:-your_username}"

echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo ""

# Confirm before proceeding
read -p "⚠️  This will modify the foreign key constraint. Continue? (y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "❌ Cancelled."
    exit 1
fi

# Run the migration
echo "📝 Executing migration..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << EOF

-- Show current constraint
SELECT 
    tc.constraint_name, 
    tc.table_name, 
    kcu.column_name,
    rc.update_rule,
    rc.delete_rule
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.referential_constraints AS rc
    ON tc.constraint_name = rc.constraint_name
WHERE tc.constraint_name = 'fk_workspace_boards';

-- Drop existing constraint
ALTER TABLE boards DROP CONSTRAINT IF EXISTS fk_workspace_boards;

-- Add new constraint with RESTRICT
ALTER TABLE boards
ADD CONSTRAINT fk_workspace_boards
FOREIGN KEY (workspace_id)
REFERENCES workspaces(id)
ON UPDATE CASCADE
ON DELETE RESTRICT;

-- Verify new constraint
SELECT 
    tc.constraint_name, 
    tc.table_name, 
    kcu.column_name,
    rc.update_rule,
    rc.delete_rule
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.referential_constraints AS rc
    ON tc.constraint_name = rc.constraint_name
WHERE tc.constraint_name = 'fk_workspace_boards';

EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Migration completed successfully!"
    echo ""
    echo "📋 Summary:"
    echo "   - Changed: boards.workspace_id foreign key"
    echo "   - From: ON DELETE CASCADE"
    echo "   - To: ON DELETE RESTRICT"
    echo ""
    echo "⚠️  Note: You can no longer delete workspaces that have boards."
    echo "   Delete all boards first, then delete the workspace."
else
    echo ""
    echo "❌ Migration failed!"
    exit 1
fi
