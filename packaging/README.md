# Packaging

## Index

- [Repository naming](#repos-naming)
- [Package naming](#package-naming)

## Topics

### Repo's naming

We suggest three ways for repository naming:

- by domain
- by purpose
- by branded name (e.g ryuk, caesar, mars)

#### By domain

Preferable services naming with using domain driven design

##### Examples

- users
- tokens
- profiles
- notifications
- lobbies

#### By purpose

Use a main purpose of application as a name

##### Examples

- users-provider
- authenticator
- profiles-manager
- notificator
- matchmaker

#### By branded name

If you have project with branded name this way is only possible to name your repo. A main area to use this way is opensource

##### Examples

- cerberus
- рeimdallr
- alfred
- bing
- cross

#### NB

We recommend not to mesh `by-domain` and `by-purpose` ways for a one project.

---

### Package naming

Package name should present certain layer or domain which is present in package.

Use singular nouns as packages' names. Using plural names is acceptable if a package presents plural entity as singular (e.g `bytes`).

##### Examples

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

- users (package presents `user` model)
- handler (package provides all server layer's functionality)
- header (package presents headers struct)

</td><td>

- user
- server
- headers

</td></tr>
</tbody></table>

---