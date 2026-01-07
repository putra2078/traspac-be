# Contacts API Documentation

## Overview
API untuk mengelola data kontak pengguna yang sedang login.

## Base URL
```
/api/v1/contacts
```

## Authentication
Semua endpoint memerlukan autentikasi menggunakan JWT token yang dikirim melalui cookie `access_token`.

---

## Endpoints

### 1. Get My Contact
Mengambil data kontak dari user yang sedang login.

**Endpoint:** `GET /api/v1/contacts/me`

**Headers:**
```
Cookie: access_token=<JWT_TOKEN>
```

**Response Success (200 OK):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "user_id": 2,
    "name": "John Doe",
    "photo": "https://your-supabase-project.supabase.co/storage/v1/object/public/your-bucket/photos/abc123.jpg",
    "email": "john.doe@example.com",
    "created_at": "2026-01-01T10:00:00Z",
    "updated_at": "2026-01-01T10:00:00Z"
  }
}
```

> **Note:** Field `photo` akan otomatis dikonversi menjadi public URL Supabase Storage jika value-nya adalah storage key (tidak dimulai dengan `http`). Ini memastikan frontend dapat langsung menggunakan URL untuk menampilkan foto.

**Response Error (404 Not Found):**
```json
{
  "status": "error",
  "message": "Contact not found"
}
```

**Response Error (401 Unauthorized):**
```json
{
  "status": "error",
  "message": "Unauthorized"
}
```

---

### 2. Update My Contact
Memperbarui data kontak dari user yang sedang login.

**Endpoint:** `PUT /api/v1/contacts/me`

**Headers:**
```
Content-Type: application/json
Cookie: access_token=<JWT_TOKEN>
```

**Request Body:**
```json
{
  "name": "John Doe Updated",
  "photo": "https://storage.example.com/photos/user123-new.jpg",
  "email": "john.updated@example.com"
}
```

**Response Success (200 OK):**
```json
{
  "status": "success",
  "message": "Contact updated successfully",
  "data": {
    "id": 1,
    "user_id": 2,
    "name": "John Doe Updated",
    "photo": "https://storage.example.com/photos/user123-new.jpg",
    "email": "john.updated@example.com",
    "created_at": "2026-01-01T10:00:00Z",
    "updated_at": "2026-01-01T12:30:00Z"
  }
}
```

**Response Error (400 Bad Request):**
```json
{
  "status": "error",
  "message": "Invalid request body"
}
```

**Response Error (401 Unauthorized):**
```json
{
  "status": "error",
  "message": "Unauthorized"
}
```

**Response Error (404 Not Found):**
```json
{
  "status": "error",
  "message": "Contact not found"
}
```

---

## Data Model

### Contact
| Field | Type | Description |
|-------|------|-------------|
| id | uint | ID unik dari contact |
| user_id | uint | ID user yang memiliki contact ini |
| name | string | Nama lengkap user |
| photo | string | URL foto profil user |
| email | string | Email user |
| created_at | timestamp | Waktu pembuatan record |
| updated_at | timestamp | Waktu update terakhir |

---

## Usage Examples

### cURL Example - Get My Contact
```bash
curl -X GET "http://localhost:8080/api/v1/contacts/me" \
  -H "Cookie: access_token=YOUR_JWT_TOKEN"
```

### cURL Example - Update My Contact
```bash
curl -X PUT "http://localhost:8080/api/v1/contacts/me" \
  -H "Content-Type: application/json" \
  -H "Cookie: access_token=YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.updated@example.com"
  }'
```

### JavaScript/Fetch Example
```javascript
// Get My Contact
async function getMyContact() {
  const response = await fetch('/api/v1/contacts/me', {
    method: 'GET',
    credentials: 'include', // Include cookies
  });
  
  const data = await response.json();
  console.log(data);
}

// Update My Contact
async function updateMyContact(contactData) {
  const response = await fetch('/api/v1/contacts/me', {
    method: 'PUT',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(contactData),
  });
  
  const data = await response.json();
  console.log(data);
}
```

---

## Notes
- Endpoint ini menggunakan `user_id` dari JWT token untuk mengambil data kontak yang sesuai
- Setiap user hanya bisa mengakses dan mengupdate data kontak miliknya sendiri
- Field `photo` otomatis dikonversi dari storage key menjadi public URL Supabase Storage
  - Jika value dimulai dengan `http`, maka akan dikembalikan apa adanya
  - Jika value adalah storage key (contoh: `photos/abc123.jpg`), akan dikonversi menjadi full public URL
- Email harus unik dalam sistem
- Photo URL yang dikembalikan dapat langsung digunakan oleh frontend untuk menampilkan gambar
