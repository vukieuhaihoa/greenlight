CREATE DATABASE greenlight;

\c greenlight

CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa5$w0rd';

CREATE EXTENSION IF NOT EXISTS citext;