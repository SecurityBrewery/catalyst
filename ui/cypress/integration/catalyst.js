describe('user', () => {

    // beforeEach(function () {
    //     cy.fixture('user_demo.json').then((user_demo) => {
    //         this.user_demo = user_demo
    //     })
    //     cy.fixture('userdata_demo.json').then((userdata_demo) => {
    //         this.userdata_demo = userdata_demo
    //     })
    // })

    it('open ticket', () => {
        cy.visit('/');

        // login
        cy.get("#username").type("bob");
        cy.get("#password").type("bob");
        cy.get("#kc-login").click();

        cy.intercept('GET', '/api/userdata/demo', { fixture: 'userdata_demo.json' })
        cy.intercept('GET', '/api/users/demo', { fixture: 'user_demo.json' })

        cy.origin('http://localhost', () => {
            cy.visit('/tickets');

            // clear caql
            cy.get("#app > div > main > div > div > div > div > header > div > div.v-input.v-input--hide-details.v-input--is-label-active.v-input--is-dirty.v-input--dense.theme--light.v-text-field.v-text-field--single-line.v-text-field--solo.v-text-field--solo-flat.v-text-field--is-booted.v-text-field--enclosed.v-text-field--placeholder > div > div > div:nth-child(2) > div > button")
                .click();

            // open ticket
            cy.get("#app > div > main > div > div > div > div > div > div.v-data-table__wrapper > table > tbody > tr:nth-child(1) > td > a")
                .click()

            // assert url
            cy.url().should('eq', "http://localhost/tickets/8123")

            // assert title
            cy.get("#\\38 123 > div > div > div:nth-child(3) > div:nth-child(2) > div:nth-child(2) > div > div.col-lg-8.col-12 > h1")
                .should("have.text", " Incident  #8123:  live zebra ")
        })
    })
})