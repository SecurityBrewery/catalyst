export interface Type {
  id: string

  singular: string
  plural: string
  icon: string
  schema: JSONSchema
}

export interface Ticket {
  id: string

  type: string
  name: string
  description: string
  open: boolean
  resolution: string
  schema: JSONSchema
  state: any
  owner: string

  created: string
  updated: string
  expand: {
    owner: User
    type: Type
    comments_via_ticket: Array<Comment>
    timeline_via_ticket: Array<TimelineItem>
    links_via_ticket: Array<Link>
    files_via_ticket: Array<File>
    tasks_via_ticket: Array<Task>
    runs_via_ticket: Array<Run>
  }
}

export interface Task {
  id: string

  ticket: string
  name: string
  open: boolean
  owner: string

  created: string
  updated: string
  expand: {
    ticket: Ticket
    owner: User
  }
}

export interface Comment {
  id: string

  ticket: string
  author: string
  message: string

  created: string
  updated: string
  expand: {
    ticket: Ticket
    author: User
  }
}

export interface TimelineItem {
  id: string

  ticket: string
  message: string
  time: string

  created: string
  updated: string
  expand: {
    ticket: Ticket
  }
}

export interface Link {
  id: string

  ticket: string
  name: string
  url: string

  created: string
  updated: string
  expand: {
    ticket: Ticket
    owner: User
  }
}

export interface File {
  id: string

  ticket: string
  name: string
  blob: string
  size: number

  created: string
  updated: string
  expand: {
    ticket: Ticket
    owner: User
  }
}

export interface User {
  id: string

  name: string
  email: string

  created: string
  updated: string
}

export interface Playbook {
  id: string

  name: string
  steps: Array<PlaybookStep>

  created: string
  updated: string
}

export interface PlaybookStep {
  name: string
  type: 'task' | 'jsonschema'
  description: string
  schema: JSONSchema
}

export interface Run {
  id: string
  ticket: string
  name: string
  steps: Array<RunStep>
}

export interface RunStep {
  name: string
  type: 'task' | 'jsonschema'
  status: 'open' | 'completed' | 'pending' | 'failed'
  description: string
  schema: JSONSchema
  state: any
}

export interface JSONSchema {
  type: 'object'
  properties: Record<
    string,
    {
      title: string
      type: string
      description?: string
    }
  >
  required?: Array<string>
}
