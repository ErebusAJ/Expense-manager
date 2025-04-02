-- +goose Up
ALTER TABLE groups 
ADD COLUMN image_url TEXT DEFAULT 'https://blogassets.airtel.in/wp-content/uploads/2023/05/felix-rostig-UmV2wr-Vbq8-unsplash.jpg';

-- +goose Down
ALTER TABLE groups
DROP COLUMN image_url;