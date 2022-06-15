// ***********************************************************
// global configuration
// https://on.cypress.io/configuration
// ***********************************************************

import './commands'

Cypress.Cookies.defaults({
    preserve: 'user',
})

Cypress.on('uncaught:exception', (err, runnable) => {
    return false
})

Cypress.Commands.add('login', (options = {}) => {
    if (Cypress.env('AUTH') === 'simple') {
        cy.contains("Name").click({force: true});
        cy.get("#username").type("tom");
        cy.contains("Password").click({force: true});
        cy.get("#password").type("tom");
        cy.get("button").contains("Login").click();
    } else if (Cypress.env('AUTH') === 'keycloak') {
        cy.get("#username").type("bob");
        cy.get("#password").type("bob");
        cy.get("#kc-login").click();
    }
})
