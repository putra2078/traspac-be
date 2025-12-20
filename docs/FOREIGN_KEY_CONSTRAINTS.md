# Foreign Key Constraints - Penjelasan

## Masalah yang Ditemukan

Ketika menghapus data `workspaces`, semua `boards` yang terkait ikut terhapus secara otomatis. Ini terjadi karena foreign key constraint menggunakan `ON DELETE CASCADE`.

## Jenis-jenis ON DELETE Actions

### 1. **CASCADE** ⚠️
```sql
ON DELETE CASCADE
```
**Behavior:** Ketika parent record dihapus, semua child records akan **ikut terhapus otomatis**.

**Contoh:**
- Hapus workspace ID 1
- Semua boards dengan workspace_id = 1 akan **otomatis terhapus**
- Semua task_tab yang terkait dengan boards tersebut juga **ikut terhapus**
- Dan seterusnya (cascade effect)

**Kapan Digunakan:**
- Ketika child data tidak berguna tanpa parent
- Contoh: `settings` table (jika user dihapus, settings-nya juga tidak berguna)

---

### 2. **RESTRICT** ✅ (Recommended untuk boards)
```sql
ON DELETE RESTRICT
```
**Behavior:** **Mencegah** penghapusan parent jika masih ada child records yang terkait.

**Contoh:**
- Coba hapus workspace ID 1
- Database akan **menolak** dengan error jika masih ada boards
- User harus menghapus semua boards terlebih dahulu

**Kapan Digunakan:**
- Ketika ingin mencegah data loss yang tidak disengaja
- Ketika child data penting dan harus dihapus secara eksplisit
- **Recommended untuk boards ↔ workspaces**

---

### 3. **SET NULL**
```sql
ON DELETE SET NULL
```
**Behavior:** Ketika parent dihapus, foreign key di child akan di-set menjadi `NULL`.

**Contoh:**
- Hapus workspace ID 1
- Semua boards dengan workspace_id = 1 akan menjadi workspace_id = NULL

**Kapan Digunakan:**
- Ketika child record masih berguna tanpa parent
- **TIDAK BISA** digunakan jika kolom foreign key adalah `NOT NULL`

---

### 4. **SET DEFAULT**
```sql
ON DELETE SET DEFAULT
```
**Behavior:** Ketika parent dihapus, foreign key di child akan di-set ke nilai default.

**Contoh:**
- Hapus workspace ID 1
- Semua boards dengan workspace_id = 1 akan menjadi workspace_id = (default value)

**Kapan Digunakan:**
- Jarang digunakan
- Ketika ada default workspace/category

---

### 5. **NO ACTION** (Default)
```sql
ON DELETE NO ACTION
```
**Behavior:** Sama seperti `RESTRICT`, tapi pengecekan dilakukan di akhir transaction.

---

## Analisis Tabel Anda

### Current Schema Analysis:

```
users (id)
  ↓ ON DELETE CASCADE
workspaces (id, created_by)
  ↓ ON DELETE CASCADE ⚠️ MASALAH DI SINI
boards (id, workspace_id)
  ↓ ON DELETE CASCADE
task_tab (id, board_id)
  ↓ ON DELETE CASCADE
task_card (id, task_tab_id)
  ↓ ON DELETE CASCADE
task_card_users (task_card_id, user_id)
```

### Rekomendasi Perubahan:

| Tabel | Foreign Key | Current | Recommended | Alasan |
|-------|-------------|---------|-------------|--------|
| **workspaces** | created_by → users | CASCADE | RESTRICT | Jangan hapus workspace jika user dihapus |
| **boards** | workspace_id → workspaces | CASCADE | **RESTRICT** | **Jangan hapus boards otomatis** |
| **task_tab** | board_id → boards | CASCADE | RESTRICT | Jangan hapus tabs otomatis |
| **task_card** | task_tab_id → task_tab | CASCADE | RESTRICT | Jangan hapus cards otomatis |
| **task_card_users** | task_card_id → task_card | CASCADE | CASCADE | OK - relasi many-to-many |
| **task_card_users** | user_id → users | CASCADE | RESTRICT | Jangan hapus assignment jika user dihapus |

## Solusi yang Diterapkan

### 1. Update Initial Migration
File: `migrations/000001_init_database_schema.up.sql`

```sql
-- BEFORE (Line 71)
ON DELETE CASCADE

-- AFTER
ON DELETE RESTRICT
```

### 2. Create New Migration
File: `migrations/000003_alter_boards_fk_constraint.up.sql`

Migration ini akan mengubah constraint yang sudah ada di database.

## Cara Menjalankan Migration

```bash
# Jalankan migration baru
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up

# Atau jika menggunakan golang-migrate
migrate -path ./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

## Testing

### Test 1: Coba Hapus Workspace yang Memiliki Boards

**BEFORE (CASCADE):**
```sql
DELETE FROM workspaces WHERE id = 1;
-- ✅ Success - boards juga ikut terhapus
```

**AFTER (RESTRICT):**
```sql
DELETE FROM workspaces WHERE id = 1;
-- ❌ ERROR: update or delete on table "workspaces" violates foreign key constraint "fk_workspace_boards" on table "boards"
-- DETAIL: Key (id)=(1) is still referenced from table "boards".
```

### Test 2: Hapus Boards Dulu, Baru Workspace

```sql
-- 1. Hapus semua boards di workspace
DELETE FROM boards WHERE workspace_id = 1;

-- 2. Baru hapus workspace
DELETE FROM workspaces WHERE id = 1;
-- ✅ Success
```

## Implementasi di Backend

Jika menggunakan RESTRICT, Anda perlu handle di backend:

```go
func (u *usecase) DeleteWorkspace(id uint) error {
    // Check if workspace has boards
    boards, err := u.boardRepo.FindByWorkspaceID(id)
    if err != nil {
        return err
    }
    
    if len(boards) > 0 {
        return errors.New("cannot delete workspace: still has boards. Please delete all boards first")
    }
    
    return u.repository.Delete(id)
}
```

Atau bisa juga dengan soft delete cascade:

```go
func (u *usecase) DeleteWorkspace(id uint) error {
    // Soft delete all boards first
    boards, _ := u.boardRepo.FindByWorkspaceID(id)
    for _, board := range boards {
        u.boardRepo.SoftDelete(board.ID)
    }
    
    // Then soft delete workspace
    return u.repository.SoftDelete(id)
}
```

## Kesimpulan

✅ **Masalah:** `ON DELETE CASCADE` menyebabkan boards terhapus otomatis ketika workspace dihapus

✅ **Solusi:** Ubah menjadi `ON DELETE RESTRICT` untuk mencegah penghapusan tidak disengaja

✅ **Benefit:** 
- Data lebih aman
- User harus eksplisit menghapus boards
- Mencegah data loss yang tidak disengaja

⚠️ **Trade-off:**
- User harus menghapus boards terlebih dahulu sebelum menghapus workspace
- Perlu handle error di backend
