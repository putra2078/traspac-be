# Panduan Integrasi Upload Foto Registrasi User

Dokumen ini menjelaskan cara mengintegrasikan fitur registrasi user baru dengan upload foto profil dari frontend ke backend.

## Endpoint

- **URL**: `/api/v1/users/`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`

## Request Body

Backend mengharapkan payload dalam format `multipart/form-data` karena menyertakan file binary.

| Field Name | Type | Required | validations | Deskripsi |
| :--- | :--- | :--- | :--- | :--- |
| `username` | Text | Yes | min 3, max 50 | Username unik pengguna |
| `email` | Text | Yes | email format | Alamat email unik |
| `password` | Text | Yes | min 8 chars | Password akun |
| `name` | Text | Yes | min 3, max 100 | Nama lengkap pengguna |
| `phone_number` | Text | Yes | min 10, max 15 | Nomor telepon |
| `gender` | Text | Yes | 'male', 'female', 'other' | Jenis kelamin |
| `birth_date` | Text | Yes | YYYY-MM-DD | Tanggal lahir |
| `address` | Text | No | - | Alamat lengkap |
| `photo_file` | File | No | - | File gambar (jpg/png) untuk foto profil |

> **Catatan:** Field `photo` (string) masih ada untuk backward compatibility jika ingin mengirim URL gambar langsung, tetapi untuk upload file gunakan `photo_file`.

## Contoh Code Integrasi (JavaScript/TypeScript)

Berikut adalah contoh fungsi menggunakan `FormData` dan `fetch` API.

```javascript
async function registerUser(userData, file) {
  const formData = new FormData();

  // Append data text
  formData.append('username', userData.username);
  formData.append('email', userData.email);
  formData.append('password', userData.password);
  formData.append('name', userData.name);
  formData.append('phone_number', userData.phoneNumber);
  formData.append('gender', userData.gender); // 'male' | 'female' | 'other'
  formData.append('birth_date', userData.birthDate); // 'YYYY-MM-DD'
  
  if (userData.address) {
    formData.append('address', userData.address);
  }

  // Append file jika ada
  // Pastikan 'file' adalah object File dari input type="file"
  if (file) {
    formData.append('photo_file', file);
  }

  try {
    const response = await fetch('http://localhost:8080/api/v1/users/', {
      method: 'POST',
      body: formData,
      // Note: Jangan set Content-Type header secara manual saat menggunakan Fetch dengan FormData.
      // Browser akan otomatis set Content-Type ke multipart/form-data dengan boundary yang benar.
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Registration failed');
    }

    const result = await response.json();
    console.log('Success:', result);
    return result;

  } catch (error) {
    console.error('Error:', error);
    throw error;
  }
}

// Penggunaan dengan input HTML
/*
<input type="file" id="profilePhoto" />
<button onclick="handleRegister()">Register</button>
*/

function handleRegister() {
    const fileInput = document.getElementById('profilePhoto');
    const file = fileInput.files[0];

    const userData = {
        username: "johndoe",
        email: "john@example.com",
        password: "secretpassword",
        name: "John Doe",
        phoneNumber: "08123456789",
        gender: "male",
        birthDate: "1990-01-01",
        address: "Jl. Sudirman No. 1"
    }

    registerUser(userData, file);
}
```

## Response

**Success 201 Created**
```json
{
  "message": "User registered successfully"
}
```

**Error 400 Bad Request** (Validasi Gagal)
```json
{
  "error": "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"
}
```

**Error 409 Conflict** (Email/Data duplicate atau error upload)
```json
{
  "error": "email already in use"
}
// atau
{
  "error": "failed to upload photo: ..."
}
```
