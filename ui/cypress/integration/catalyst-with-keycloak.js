describe('user', () => {
    it('open ticket', () => {
        cy.visit('/');

        // login
        cy.get("#username").type("bob");
        cy.get("#password").type("bob");
        cy.get("#kc-login").click();

        cy.intercept('GET', '/api/userdata/demo', { fixture: 'userdata_demo.json' })
        cy.intercept('GET', '/api/users/demo', { fixture: 'user_demo.json' })

        cy.origin('http://localhost/ui/', () => {
            cy.visit('/tickets');

            // clear caql
            cy.get("#caqlbar > div > div > div > div > div:nth-child(2) > div > button").click();

            // open ticket
            cy.contains("live zebra").click()

            // assert url
            cy.url().should('eq', "http://localhost/ui/tickets/8123")

            // assert title
            cy.get("h1").should("have.text", " Incident  #8123:  live zebra ")
        })
    })
})