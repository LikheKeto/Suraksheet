ALTER TABLE bins 
ADD CONSTRAINT uq_bin_name 
UNIQUE (name, owner);