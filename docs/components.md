### Custom Components
This page will document some of the custom components of transportd and how they can be used by consumers.

##### Validate Headers

The `validate_headers` component can be used to validate the presence of headers and the existence of specific values in those headers.
It verifies that *any* of the specified headers exist with *any* of the specified values for that header. Each key under `allowed`
is the header name and the list of values for that key represent the the values we want to verify exist.

If a configured header does not exist, or the specified value is not found, it will check for the next `allowed` header, if configured,
and if the configured value exists. If it can't find any configured header or configured value for that header, it will reject the request.
This makes the plugin flexible to check for the presence of any header values combination. It does not currently do strict matching
if you want to validate that *multiple* values exist.

Example: This configuration would verify that a header named `LdapGroups` exists and that it has any of the listed values `hr` or `pm`.
```yaml
validateheaders:
  allowed:
    ldap-groups:
      - "hr"
      - "pm"
```


Example: This configuration would verify that a header named `LdapGroups` or `Username` exists with any of the listed values for that header.
```yaml
validateheaders:
  allowed:
    ldap-groups:
      - "hr"
      - "pm"
    username:
      - "cloud-admin"
```
