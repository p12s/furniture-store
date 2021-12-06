# Furniture store

Technical task is [here](tech_ru.md)  
System Design - microservices with asynchronous interaction is [here](ddd_ru.md)  

## Event storming schema
![ES](https://github.com/p12s/furniture-store/blob/master/images/ES.png?raw=true)  
[link](https://lucid.app/lucidchart/1482e706-4b6d-49f8-adce-e0b7932d8bbe/edit?viewport_loc=-128%2C-54%2C2307%2C1397%2C0_0&invitationId=inv_dd15d087-fe4e-4cea-b2f5-ce0f5ad99f35)  

## Domain model
![Domain model](https://github.com/p12s/furniture-store/blob/master/images/domain-model.png?raw=true)  
[link](https://www.xmind.net/m/EVD7bc)  

### Disclaimer  
If you have any comments or think that you have found a mistake - feel free to create an issue!  







internal/repository/repository.go:44:3: exitAfterDefer: log.Fatal will exit, and `defer statement.Close()` will not run (gocritic)
		log.Fatal("create account.account table fail: ", err.Error())
		^
internal/transport/rest/account.go:10:19: func `(*Handler).updateAccount` is unused (unused)
func (h *Handler) updateAccount(c *gin.Context) {
                  ^
internal/transport/rest/middleware.go:48:6: func `getAccountId` is unused (unused)
func getAccountId(c *gin.Context) (int, error) {
     ^
internal/transport/rest/auth.go:58:19: func `(*Handler).token` is unused (unused)
func (h *Handler) token(c *gin.Context) {