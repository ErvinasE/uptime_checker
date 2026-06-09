-- Use endpoints that respond reliably to automated checks.
UPDATE websites SET url = 'https://registry.npmjs.org' WHERE name = 'npm';
