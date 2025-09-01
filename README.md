<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<!-- <a id="readme-top"></a>
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url] -->
<!-- [![project_license][license-shield]][license-url] -->
<!-- [![LinkedIn][linkedin-shield]][linkedin-url] -->

<!-- PROJECT LOGO -->
<br />
<div align="center">
<h3 align="center">SIMEDIS REST - Sistem Informasi Rekam Medis Rest Api</h3>
</div>
<br>

<!-- ABOUT THE PROJECT -->
## About
Ini adalah backend REST API sederhana untuk **SIMEDIS (Sistem Informasi Rekam Medis)**, sebuah aplikasi yang dirancang untuk mengelola data klinis di fasilitas kesehatan seperti puskesmas. Proyek ini merupakan hasil re-engineering dari aplikasi sebelumnya berbasis laravel (https://github.com/franklindh/simedis) menjadi sebuah sistem modern menggunakan Go

![Go](https://img.shields.io/badge/Go-white?style=for-the-badge&logo=go&logoColor=blue)
![Gin](https://img.shields.io/badge/Gin-white?style=for-the-badge&logo=gin&logoColor=blue)
![Database](https://img.shields.io/badge/PostgreSQL-white?style=for-the-badge&logo=postgresql)

<!-- <p align="right">(<a href="#readme-top">back to top</a>)</p> -->
## Feature
* **Manajemen Petugas**: CRUD untuk data petugas (Admin, Dokter, Poli, Lab) dengan sistem *role-based*.
* **Manajemen Pasien**: CRUD untuk data demografi dan rekam medis pasien.
* **Manajemen Master Data**: Pengelolaan data poliklinik, jadwal dokter, dan klasifikasi penyakit (ICD).
* **Alur Klinis**:
    * Pendaftaran antrian pasien ke jadwal dokter yang tersedia.
    * Pembuatan rekam medis (pemeriksaan) yang terhubung ke data antrian.
    * Pencatatan hasil laboratorium.
* **Laporan**: Agregasi data untuk laporan kunjungan dan penyakit terbanyak.

<!-- GETTING STARTED -->

## Setup
### Instalasi

1. Clone repo
   ```sh
   git clone https://github.com/franklindh/simedis-rest.git
   ```

2. Copy file `.env.example` menjadi `.env`
   ```bash
   cp .env.example .env
2. Install dependensi
   ```sh
   go mod tidy
    ```
3. Jalankan 
   ```sh
   go run cmd/api/server.go
   ```
<!-- <p align="right">(<a href="#readme-top">back to top</a>)</p> -->



<!-- USAGE EXAMPLES -->
<!-- ## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_ -->

<!-- <p align="right">(<a href="#readme-top">back to top</a>)</p> -->

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/franklindh/simedis-rest.svg?style=for-the-badge
[contributors-url]: https://github.com/franklindh/simedis-rest/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/franklindh/simedis-rest.svg?style=for-the-badge
[forks-url]: https://github.com/franklindh/simedis-rest/network/members
[stars-shield]: https://img.shields.io/github/stars/franklindh/simedis-rest.svg?style=for-the-badge
[stars-url]: https://github.com/franklindh/simedis-rest/stargazers
[issues-shield]: https://img.shields.io/github/issues/franklindh/simedis-rest.svg?style=for-the-badge
[issues-url]: https://github.com/franklindh/simedis-rest/issues

[Go]: https://img.shields.io/badge/go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[Gin]: https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=gin&logoColor=white
[Gin-url]: https://gin-gonic.com/
[Postgresql]: https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql&logoColor=blue
[Postgresql-url]: https://www.postgresql.org/
[Angular.io]: https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white
[Angular-url]: https://angular.io/
[Svelte.dev]: https://img.shields.io/badge/Svelte-4A4A55?style=for-the-badge&logo=svelte&logoColor=FF3E00
[Svelte-url]: https://svelte.dev/
[Laravel.com]: https://img.shields.io/badge/Laravel-FF2D20?style=for-the-badge&logo=laravel&logoColor=white
[Laravel-url]: https://laravel.com
[Bootstrap.com]: https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white
[Bootstrap-url]: https://getbootstrap.com
[JQuery.com]: https://img.shields.io/badge/jQuery-0769AD?style=for-the-badge&logo=jquery&logoColor=white
[JQuery-url]: https://jquery.com 