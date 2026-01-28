package authz

import future.keywords.contains
import future.keywords.if
import future.keywords.in
import future.keywords.every

abilities contains ability if {
  some definition in data.role_definitions
  
  # Check if all grants in the definition are present in the input
  every grant in definition.grants {
    grant in input.grants
  }
  
  ability := definition.ability
}
