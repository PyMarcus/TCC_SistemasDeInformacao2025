ALTER TABLE atoms DROP CONSTRAINT IF EXISTS atoms_datasets_id_fkey;
ALTER TABLE atoms DROP CONSTRAINT IF EXISTS atoms_question_id_fkey;
ALTER TABLE atoms DROP CONSTRAINT IF EXISTS atoms_error_id_fkey;

DROP TABLE IF EXISTS atoms;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS errors;
DROP TABLE IF EXISTS datasets;