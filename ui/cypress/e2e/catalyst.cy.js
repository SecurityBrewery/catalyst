beforeEach(() => {
    cy.visit('/');

    cy.login();

    cy.getCookie('user').should('exist');

    cy.intercept('GET', '/api/userdata/demo', { fixture: 'userdata_demo.json' })
    cy.intercept('GET', '/api/users/demo', { fixture: 'user_demo.json' })
})

if (Cypress.env('TEST') === 'tickets') {
    describe('tickets', () => {
        it('open ticket', () => {
            cy.visit('http://localhost/ui/tickets');

            // clear caql
            cy.get("#caqlbar > div > div > div > div > div:nth-child(2) > div > button").click();

            // wait for results to load
            cy.get("#app > div > main > div > div > div > div > div > div.v-data-table__wrapper > table > tbody > tr:nth-child(1)").should('exist');
            cy.wait(1000);

            // open ticket
            cy.contains("live zebra").click();

            // assert url
            cy.url().should('eq', "http://localhost/ui/tickets/8123")

            // assert title
            cy.get("h1").should("have.text", " Incident  #8123:  live zebra ")
        })
    })
}

if (Cypress.env('TEST') === 'templates') {
    describe('templates', () => {
        it('create template', () => {
            cy.get("#toggle_menu").click();
            cy.contains('Templates').click();
            cy.get("#toggle_menu").click();

            cy.get("body").then($body => {
                if ($body.find('a[href="/ui/templates/description-only"]').length > 0) {
                    cy.get('a[href="/ui/templates/description-only"]').trigger('mouseover');
                    cy.get('a[href="/ui/templates/description-only"] button').click();
                    cy.get("#delete-button").click();
                }
            });

            cy.contains("New Template").click();
            cy.url().should('eq', "http://localhost/ui/templates/new");
            cy.get("#name-edit").click().clear().type("DescriptionOnly");
            cy.get("#template-edit #advanced").click({force: true});
            cy.get(".prism-editor__textarea").clear().type('{ "type": "object", "name": "Incident", "required": [ "description" ], "properties": { "description": { "title": "Description", "type": "string", "x-display": "textarea" } } }', {parseSpecialCharSequences: false});
            cy.contains("Create").click();

            cy.get("#toggle_menu").click();
            cy.contains('Alerts').click();
            cy.get("#toggle_menu").click();

            cy.contains("New Alert").click();
            cy.get("#title-edit").type("New Alert");
            cy.get("#templates-edit .v-input__slot").click();
            cy.contains("DescriptionOnly").click();
            cy.contains("Create").click();
            cy.url().should('match', /http:\/\/localhost\/ui\/tickets\/\d+/);

            cy.get("#description").type("Lorem Ipsum");
            cy.contains("Save Details").click();
        })
    })
}

if (Cypress.env('TEST') === 'playbooks') {
    describe('playbooks', () => {
        it('create playbook', () => {
            cy.get("#toggle_menu").click();
            cy.contains('Playbooks').click();
            cy.get("#toggle_menu").click();

            cy.get("body").then($body => {
                if ($body.find('a[href="/ui/playbooks/test"]').length > 0) {
                    cy.get('a[href="/ui/playbooks/test"]').trigger('mouseover');
                    cy.get('a[href="/ui/playbooks/test"] button').click();
                    cy.get("#delete-button").click();
                }
            });

            cy.contains("New Playbook").click();
            cy.url().should('eq', "http://localhost/ui/playbooks/new");
            cy.get(".prism-editor__textarea").clear().type('name: Test\n' +
                'tasks:\n' +
                '  input:\n' +
                '  name: Enter something to hash\n' +
                'type: input\n' +
                'schema:\n' +
                '  title: Something\n' +
                'type: object\n' +
                'properties:\n' +
                '  something:\n' +
                '  type: string\n' +
                'title: Something\n' +
                'default: ""\n' +
                '{backspace}{backspace}{backspace}next:\n' +
                '  hash: "something != \'\'"\n' +
                '{backspace}{backspace}\n' +
                'hash:\n' +
                '  name: Hash the something\n' +
                'type: automation\n' +
                'automation: hash.sha1\n' +
                'payload:\n' +
                '  default: "playbook.tasks[\'input\'].data[\'something\']"\n' +
                '{backspace}next:\n' +
                '  comment: "hash != \'\'"\n' +
                '{backspace}{backspace}\n' +
                'comment:\n' +
                '  name: Comment the hash\n' +
                'type: automation\n' +
                'automation: comment\n' +
                'payload:\n' +
                '  default: "playbook.tasks[\'hash\'].data[\'hash\']"\n' +
                '{backspace}next:\n' +
                '  done: "done"\n' +
                '{backspace}{backspace}\n' +
                'done:\n' +
                '  name: You can close this case now\n' +
                'type: task\n');
            cy.scrollTo('bottom');
            cy.contains("Create").click();

            cy.get("#toggle_menu").click();
            cy.contains('Alerts').click();
            cy.get("#toggle_menu").click();

            cy.contains("New Alert").click();
            cy.get("#title-edit").type("New Alert");
            cy.get("#playbooks-edit .v-input__slot").click();
            cy.contains("Test").click();
            cy.contains("Create").click();
            cy.url().should('match', /http:\/\/localhost\/ui\/tickets\/\d+/);

            cy.scrollTo('bottom');
            cy.get(".playbook-test .tasks").contains("Enter something to hash").click();
            cy.get("#something").type("my test value");
            cy.contains("Complete").click();

            // cy.wait(300 * 1000);
            // cy.contains("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3").should('exist');
        })
    })
}
