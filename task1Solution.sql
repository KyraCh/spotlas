CREATE TEMPORARY TABLE TempTable AS
SELECT 
    id,
    name,
    CASE
        WHEN SUBSTRING_INDEX(website, '://', -1) LIKE '%bit.ly%' OR 
             SUBSTRING_INDEX(website, '://', -1) LIKE '%goo.gl%' OR 
             SUBSTRING_INDEX(website, '://', -1) LIKE '%fb.me%' OR 
             SUBSTRING_INDEX(website, '://', -1) LIKE '%t.co%' OR 
             SUBSTRING_INDEX(website, '://', -1) LIKE '%ow.ly%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%tinyurl.com%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%youtu.be%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%linktr.ee%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%github.com%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%linkedin.com%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%facebook.com%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%twitter.com%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%instagram.com%' OR
             SUBSTRING_INDEX(website, '://', -1) LIKE '%youtube.com%' THEN
            SUBSTRING_INDEX(SUBSTRING_INDEX(website, '://', -1), 'www.', -1)
        ELSE
            SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(website, '://', -1), '/', 1), 'www.', -1)
    END as domain

FROM 
    MY_TABLE;


CREATE TEMPORARY TABLE TempTableCount AS
SELECT 
    domain, 
    COUNT(*) as domain_count
FROM 
    TempTable 
GROUP BY 
    domain;

SELECT 
    t1.name, 
    t1.domain, 
    t2.domain_count 
FROM 
    TempTable t1 
JOIN 
    TempTableCount t2 
ON 
    t1.domain = t2.domain
WHERE 
    t2.domain_count > 1;
