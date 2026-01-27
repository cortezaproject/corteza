import { Event } from '../eventbus/shared'
import { Role } from './types/role'
import { User } from './types/user'
interface TriggerEndpoints {
    automationTriggerScript(params: {
        script: string;
    }): Promise<object>;
    roleTriggerScript(params: {
        roleID: string;
        script: string;
    }): Promise<object>;
    userTriggerScript(params: {
        userID: string;
        script: string;
    }): Promise<object>;
}
export declare function SystemEvent(eventType?: string): Event;
export declare function UserEvent(user: User, eventType?: string): Event;
export declare function RoleEvent(role: Role, eventType?: string): Event;
/**
 * Returns handler that routes onManual events for server script to the system API
 *
 * See makeAutomationScriptsRegistrator
 *
 * @param api
 * @return function
 */
export declare function TriggerSystemServerScriptOnManual(api: TriggerEndpoints): (ev: Event, script: string) => Promise<unknown>;
export {}
