
# Array of strings containing the id of the users
output "users_list" {
  value = [
    conduktor_console_user_v2.user1.name,
    conduktor_console_user_v2.user2.name
  ]
}

