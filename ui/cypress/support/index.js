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
    cy.get("#app > div > div > form > div > div.v-card__text.text-center > div:nth-child(1) > div > div.v-input__slot > div > label")
        .click({force: true});
    cy.get("#input-7").type("tom");
    cy.get("#app > div > div > form > div > div.v-card__text.text-center > div:nth-child(2) > div > div.v-input__slot > div > label")
        .click({force: true});
    cy.get("#input-10").type("tom");
    cy.get("#app > div > div > form > div > div.v-card__actions > button").click();

    cy.getCookie('user').should('exist');
})
