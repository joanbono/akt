![](img/AKT_Banner.svg)

## USAGE

Rotate Access Keys for the default user on `default` profile. The keys will be printed and the user has to save them:

```sh
$ akt -rotate

[default]
aws_access_key_id = ${NEW_ACCESS_KEY}
aws_secret_access_key = ${NEW_SECRET_KEY}
```

Rotate Access Keys for the user `joanbono` on `account` profile, and save the new keys into `.aws/credentials`:

```sh
$ akt -user joanbono -profile account -save
```

Diffing the `.aws/credentials` file, the new key is saved there:

```diff
diff --git a/credentials.backup b/credentials
index 1123ef2..83a856e 100644
--- a/credentials.backup
+++ b/credentials
@@ -33,2 +33,2 @@ 
-aws_access_key_id = ${OLD_ACCESS_KEY}
-aws_secret_access_key = ${OLD_SECRET_KEY}
+aws_access_key_id = ${NEW_ACCESS_KEY}
+aws_secret_access_key = ${NEW_SECRET_KEY}
```

## TODO

+ [ ] Bulk key rotation. List of usernames. Generate JSON/CSV output with the new generated keys