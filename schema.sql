DROP TABLE IF EXISTS aksi_sosial CASCADE;
DROP TABLE IF EXISTS transaksi_jasa CASCADE;
DROP TABLE IF EXISTS iklan_jasa CASCADE;
DROP TABLE IF EXISTS siswa CASCADE;

CREATE TABLE siswa (
  id          SERIAL PRIMARY KEY,
  nama        VARCHAR(100) NOT NULL,
  email       VARCHAR(150) UNIQUE NOT NULL,
  password    VARCHAR(255) NOT NULL,
  kelas       VARCHAR(50)  NOT NULL,
  jurusan     VARCHAR(100) NOT NULL,
  saldo_poin  INTEGER      DEFAULT 7,
  created_at  TIMESTAMP    DEFAULT NOW()
);

CREATE TABLE iklan_jasa (
  id           SERIAL PRIMARY KEY,
  penyedia_id  INTEGER REFERENCES siswa(id) ON DELETE CASCADE,
  judul        VARCHAR(150) NOT NULL,
  kategori     VARCHAR(50)  NOT NULL,
  deskripsi    TEXT         NOT NULL,
  harga_poin   INTEGER      NOT NULL,
  status       VARCHAR(20)  DEFAULT 'tersedia',
  created_at   TIMESTAMP    DEFAULT NOW()
);

CREATE TABLE transaksi_jasa (
  id           SERIAL PRIMARY KEY,
  iklan_id     INTEGER REFERENCES iklan_jasa(id) ON DELETE CASCADE,
  pembeli_id   INTEGER REFERENCES siswa(id) ON DELETE CASCADE,
  penyedia_id  INTEGER REFERENCES siswa(id) ON DELETE CASCADE,
  jumlah_poin  INTEGER      NOT NULL,
  status       VARCHAR(20)  DEFAULT 'berjalan',
  created_at   TIMESTAMP    DEFAULT NOW()
);

CREATE TABLE aksi_sosial (
  id              SERIAL PRIMARY KEY,
  siswa_id        INTEGER REFERENCES siswa(id) ON DELETE CASCADE,
  deskripsi       TEXT         NOT NULL,
  foto_url        VARCHAR(255),
  status          VARCHAR(20)  DEFAULT 'menunggu',
  poin_diberikan  INTEGER      DEFAULT 5,
  created_at      TIMESTAMP    DEFAULT NOW()
);
