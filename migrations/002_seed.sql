
TRUNCATE TABLE seats, seat_locks, tickets, refunds, transactions, schedules, users, movies, studios, cinemas CASCADE;

INSERT INTO cinemas (id, name, city, address) VALUES
  ('10000000-0000-0000-0000-000000000001', 'CinemaXX Grand Indonesia', 'Jakarta', 'Jl. M.H. Thamrin No.1, Jakarta Pusat'),
  ('10000000-0000-0000-0000-000000000002', 'CinemaXX Bandung Indah Plaza', 'Bandung', 'Jl. Merdeka No.56, Bandung'),
  ('10000000-0000-0000-0000-000000000003', 'CinemaXX Malioboro Mall', 'Yogyakarta', 'Jl. Malioboro No.52-58, Yogyakarta');


INSERT INTO studios (id, cinema_id, name, capacity) VALUES
  ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'Studio 1', 100),
  ('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', 'Studio 2 IMAX', 80),
  ('20000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000002', 'Studio 1', 90),
  ('20000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000003', 'Studio 1', 75);


DO $$
DECLARE
  rows CHAR(1)[] := ARRAY['A','B','C','D','E','F','G','H','I','J'];
  r CHAR(1);
  n INT;
BEGIN
  FOREACH r IN ARRAY rows LOOP
    FOR n IN 1..10 LOOP
      INSERT INTO seats (id, studio_id, row_label, seat_number, seat_type)
      VALUES (
        gen_random_uuid(),
        '20000000-0000-0000-0000-000000000001',
        r, n,
        CASE WHEN r IN ('A','B') THEN 'vip'::seat_type
             WHEN r IN ('C','D') THEN 'premium'::seat_type
             ELSE 'regular'::seat_type END
      );
    END LOOP;
  END LOOP;
END $$;


DO $$
DECLARE
  rows CHAR(1)[] := ARRAY['A','B','C','D','E','F','G','H'];
  r CHAR(1);
  n INT;
BEGIN
  FOREACH r IN ARRAY rows LOOP
    FOR n IN 1..10 LOOP
      INSERT INTO seats (id, studio_id, row_label, seat_number, seat_type)
      VALUES (
        gen_random_uuid(),
        '20000000-0000-0000-0000-000000000002',
        r, n,
        CASE WHEN r IN ('A','B') THEN 'vip'::seat_type
             WHEN r IN ('C','D') THEN 'premium'::seat_type
             ELSE 'regular'::seat_type END
      );
    END LOOP;
  END LOOP;
END $$;


INSERT INTO movies (id, title, duration_minutes, genre, rating, synopsis, poster_url) VALUES
  (
    '30000000-0000-0000-0000-000000000001',
    'Laskar Pelangi 2',
    120,
    'Drama',
    'SU',
    'Sekuel film ikonik Laskar Pelangi yang mengisahkan perjalanan anak-anak Belitung meraih impian mereka.',
    'https://example.com/posters/laskar-pelangi-2.jpg'
  ),
  (
    '30000000-0000-0000-0000-000000000002',
    'Gundala Returns',
    135,
    'Action / Superhero',
    '13+',
    'Sancaka kembali sebagai Gundala untuk menghadapi ancaman baru yang mengancam keamanan Jakarta.',
    'https://example.com/posters/gundala-returns.jpg'
  ),
  (
    '30000000-0000-0000-0000-000000000003',
    'Horor Pesantren',
    105,
    'Horror',
    '17+',
    'Sekelompok santri menemukan misteri gelap di balik tembok pesantren tua yang terisolasi.',
    'https://example.com/posters/horor-pesantren.jpg'
  );


INSERT INTO schedules (id, movie_id, studio_id, show_time, price, status) VALUES
  (
    '40000000-0000-0000-0000-000000000001',
    '30000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    NOW() + INTERVAL '1 day' + INTERVAL '10 hours',
    50000,
    'active'
  ),
  (
    '40000000-0000-0000-0000-000000000002',
    '30000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    NOW() + INTERVAL '1 day' + INTERVAL '13 hours',
    50000,
    'active'
  ),
  (
    '40000000-0000-0000-0000-000000000003',
    '30000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000002',
    NOW() + INTERVAL '1 day' + INTERVAL '14 hours',
    85000,
    'active'
  ),
  (
    '40000000-0000-0000-0000-000000000004',
    '30000000-0000-0000-0000-000000000003',
    '20000000-0000-0000-0000-000000000001',
    NOW() + INTERVAL '2 days' + INTERVAL '19 hours',
    60000,
    'active'
  ),
  (
    '40000000-0000-0000-0000-000000000005',
    '30000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000003',
    NOW() + INTERVAL '1 day' + INTERVAL '16 hours',
    75000,
    'active'
  );


INSERT INTO users (id, name, email, password_hash, phone, role) VALUES
  (
    '50000000-0000-0000-0000-000000000001',
    'Admin MKP',
    'admin@mkp.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    '081200000001',
    'admin'
  );