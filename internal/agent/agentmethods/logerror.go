package agentmethods

import (

	"fmt"
	"log"
)



// Логирование ошибки
func (a *Agent) logError(msg string, err error) {
    log.Printf("%s: %v\n", msg, err)
}

// Обработка ошибок и пропуск шага
func (a *Agent) handleErrorAndContinue(action string, err error) {
    a.logError(fmt.Sprintf("Ошибка %s", action), err)
}