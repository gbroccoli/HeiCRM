package email

import "fmt"

// RenderWelcome generates welcome email HTML with credentials
func RenderWelcome(name, emailAddr, password string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <h2 style="color: #2c3e50;">Добро пожаловать в HeiCRM!</h2>
  <p>Здравствуйте, <strong>%s</strong>!</p>
  <p>Ваша учётная запись успешно создана. Ниже указаны данные для входа:</p>
  <div style="background: #f8f9fa; padding: 15px; border-radius: 5px; margin: 15px 0;">
    <p style="margin: 5px 0;"><strong>Email:</strong> %s</p>
    <p style="margin: 5px 0;"><strong>Пароль:</strong> %s</p>
  </div>
  <p style="color: #e74c3c;"><strong>Важно:</strong> Рекомендуем сменить пароль после первого входа.</p>
  <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
  <p style="color: #999; font-size: 12px;">HeiCRM — система управления общежитием</p>
</body>
</html>`, name, emailAddr, password)
}

// RenderTaskAssigned generates task assignment notification email
func RenderTaskAssigned(taskType, priority, description string) string {
	priorityColor := map[string]string{
		"low": "#27ae60", "medium": "#f39c12", "high": "#e67e22", "critical": "#e74c3c",
	}
	color := priorityColor[priority]
	if color == "" {
		color = "#333"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <h2 style="color: #2c3e50;">Вам назначена заявка</h2>
  <div style="background: #f8f9fa; padding: 15px; border-radius: 5px; margin: 15px 0;">
    <p style="margin: 5px 0;"><strong>Тип:</strong> %s</p>
    <p style="margin: 5px 0;"><strong>Приоритет:</strong> <span style="color: %s;">%s</span></p>
    <p style="margin: 5px 0;"><strong>Описание:</strong> %s</p>
  </div>
  <p>Пожалуйста, ознакомьтесь с заявкой в системе HeiCRM.</p>
  <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
  <p style="color: #999; font-size: 12px;">HeiCRM — система управления общежитием</p>
</body>
</html>`, taskType, color, priority, description)
}

// RenderTaskStatusChanged generates task status change notification email
func RenderTaskStatusChanged(taskType, prevStatus, newStatus, description string) string {
	statusNames := map[string]string{
		"new": "Новая", "assigned": "Назначена", "in_progress": "В работе",
		"completed": "Выполнена", "closed": "Закрыта",
	}
	prev := statusNames[prevStatus]
	if prev == "" {
		prev = prevStatus
	}
	next := statusNames[newStatus]
	if next == "" {
		next = newStatus
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <h2 style="color: #2c3e50;">Статус заявки изменён</h2>
  <div style="background: #f8f9fa; padding: 15px; border-radius: 5px; margin: 15px 0;">
    <p style="margin: 5px 0;"><strong>Тип:</strong> %s</p>
    <p style="margin: 5px 0;"><strong>Статус:</strong> %s → <strong>%s</strong></p>
    <p style="margin: 5px 0;"><strong>Описание:</strong> %s</p>
  </div>
  <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
  <p style="color: #999; font-size: 12px;">HeiCRM — система управления общежитием</p>
</body>
</html>`, taskType, prev, next, description)
}
