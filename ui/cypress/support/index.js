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
    cy.visit('/');

    // login
    cy.contains("Name").click({force: true});
    cy.get("#username").type("tom");
    cy.contains("Password").click({force: true});
    cy.get("#password").type("tom");
    cy.get("button").contains("Login").click();

    cy.getCookie('user').should('exist');
})
