describe('Login', () => {
    it('can login', () => {
        cy.visit('/home');

        // login
        cy.get("#username").type("bob");
        cy.get("#password").type("bob");
        cy.get("#kc-login").click();

        cy.intercept('GET', '/api/userdata/demo', {
            statusCode: 200,
            body: {"email": "demo@example.org", "id": "demo", "name": "Demo"},
        })
        cy.intercept('GET', '/api/users/demo', {
            statusCode: 200,
            body: {
                "apikey": false,
                "blocked": false,
                "id": "demo",
                "roles": ["admin:backup:read", "admin:backup:restore", "admin:dashboard:write", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:settings:write", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"]
            },
        })

        cy.origin('http://localhost', () => {
            cy.visit('/tickets');

            // clear caql
            cy.get("#app > div > main > div > div > div > div > header > div > div.v-input.v-input--hide-details.v-input--is-label-active.v-input--is-dirty.v-input--dense.theme--light.v-text-field.v-text-field--single-line.v-text-field--solo.v-text-field--solo-flat.v-text-field--is-booted.v-text-field--enclosed.v-text-field--placeholder > div > div > div:nth-child(2) > div > button")
                .click();

            // open ticket
            cy.get("#app > div > main > div > div > div > div > div > div.v-data-table__wrapper > table > tbody > tr:nth-child(1) > td > a")
                .click()

            cy.url().should('eq', "http://localhost/tickets/8123")

            cy.get("#\\38 123 > div > div > div:nth-child(3) > div:nth-child(2) > div:nth-child(2) > div > div.col-lg-8.col-12 > h1")
                .should("have.text", " Incident  #8123:  live zebra ")
        })
    })
})