// ***********************************************************
// global configuration
// https://on.cypress.io/configuration
// ***********************************************************

import './commands'

Cypress.Cookies.defaults({
    preserve: 'user',
})

Cypress.Commands.add('login', (options = {}) => {
    cy.visit('/');

    // login
    cy.get("body > div > div > div > form > div:nth-child(1) > label").click();
    cy.get("#username").type("tom");
    cy.get("body > div > div > div > form > div:nth-child(2) > label").click();
    cy.get("#password").type("tom");
    cy.get("#submit").click();

    cy.getCookie('user').should('exist');
})